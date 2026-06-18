package goal

import (
	"context"
	"controle_financas/internal/core/cache"
	"controle_financas/internal/core/domain/apiError"
	"controle_financas/internal/core/filters"
	"controle_financas/internal/core/transaction"
	"controle_financas/internal/core/validator"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo       goalRepository
	tx         transaction.Manager
	cache      cache.Cache
	keyBuilder cache.KeyBuilder
}

func NewService(
	repo goalRepository,
	tx transaction.Manager,
	cache cache.Cache,
	keyBuilder cache.KeyBuilder,
) *Service {
	return &Service{
		repo:       repo,
		tx:         tx,
		cache:      cache,
		keyBuilder: keyBuilder,
	}
}

func (s *Service) FindAll(
	ctx context.Context,
	name string,
	f filters.Filters,
) ([]*Goal, filters.Metadata, error) {
	key := s.keyBuilder.BuildListKey(f.Page, f.Sort, f.PageSize, f.Limit(), f.Offset(), name)

	type cachePayload struct {
		Data     []*Goal
		Metadata filters.Metadata
	}

	payload, err := cache.FetchOrCache(ctx, s.cache, key, func() (cachePayload, error) {
		models, metadata, err := s.repo.FindAllByUserId(ctx, name, f)

		for _, g := range models {
			if s.verifyFailedStatus(g) {
				g.Status = GoalStatusFailed
			}
			s.calculateInstallments(g)
		}

		return cachePayload{Data: models, Metadata: metadata}, err
	})

	if err != nil {
		return nil, filters.Metadata{}, err
	}

	return payload.Data, payload.Metadata, nil
}

func (s *Service) FindById(
	ctx context.Context,
	id uuid.UUID,
) (*Goal, error) {
	key := s.keyBuilder.BuildItemKey(id.String())

	return cache.FetchOrCache(ctx, s.cache, key, func() (*Goal, error) {
		g, err := s.repo.FindById(ctx, id)
		if err != nil {
			return nil, err
		}

		err = s.tx.RunInTx(ctx, func(ctx context.Context) error {
			return s.updateStatus(ctx, g)
		})

		if err != nil {
			return nil, err
		}

		s.calculateInstallments(g)

		return g, nil
	})
}

func (s *Service) Insert(
	ctx context.Context,
	model *Goal,
) error {
	v := validator.New()
	if model.ValidateGoal(v); !v.Valid() {
		return apiError.NewValidationError(v.Errors)
	}

	err := s.tx.RunInTx(ctx, func(ctx context.Context) error {
		return s.repo.Insert(ctx, model)
	})

	if err != nil {
		return err
	}

	go func() {
		_ = s.cache.DeleteByPrefix(context.Background(), s.keyBuilder.GetPrefix())
	}()

	return nil
}

func (s *Service) Update(
	ctx context.Context,
	model *Goal,
) error {
	v := validator.New()
	if model.ValidateGoal(v); !v.Valid() {
		return apiError.NewValidationError(v.Errors)
	}
	err := s.tx.RunInTx(ctx, func(ctx context.Context) error {
		return s.repo.Update(ctx, model)
	})

	if err != nil {
		return err
	}

	go func() {
		_ = s.cache.DeleteByPrefix(context.Background(), s.keyBuilder.GetPrefix())
	}()
	return nil
}

func (s *Service) updateStatus(
	ctx context.Context,
	g *Goal,
) error {
	if s.verifyFailedStatus(g) {
		g.Status = GoalStatusFailed
		return s.tx.RunInTx(ctx, func(ctx context.Context) error {
			return s.Update(ctx, g)
		})
	}
	return nil
}

func (s *Service) DeleteById(
	ctx context.Context,
	id uuid.UUID,
) error {
	err := s.tx.RunInTx(ctx, func(ctx context.Context) error {
		return s.repo.DeleteById(ctx, id)
	})

	if err != nil {
		return err
	}

	go func() {
		_ = s.cache.DeleteByPrefix(context.Background(), s.keyBuilder.GetPrefix())
	}()

	return nil
}

func (s *Service) verifyFailedStatus(goal *Goal) bool {
	return time.Now().After(goal.Deadline) && goal.Current < goal.Amount && goal.Status != GoalStatusFailed
}

func (s *Service) calculateInstallments(g *Goal) {
	now := time.Now()
	if g.Deadline.Before(now) {
		return
	}

	if g.Installments == nil {
		g.Installments = &Installments{}
	}

	yearDiff := g.Deadline.Year() - now.Year()
	monthDiff := int(g.Deadline.Month()) - int(now.Month())
	quantity := yearDiff*12 + monthDiff

	if quantity <= 0 {
		return
	}

	g.Installments.Quantity = quantity
	remaing := g.Amount - g.Current
	if remaing <= 0 {
		return
	}

	g.Installments.Amount = remaing / float64(quantity)
}
