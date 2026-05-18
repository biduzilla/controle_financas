package transaction

import (
	"fmt"
	"net/http"
	"time"

	"controle_financas/internal/core/domain/apiError"
	"controle_financas/internal/core/handler"
	"controle_financas/utils"

	"github.com/google/uuid"
)

type transactionHandler struct {
	service    TransactionService
	errHandler apiError.ErrorHandler
}

type TransactionHandler interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	FindByID(w http.ResponseWriter, r *http.Request)
	Save(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	DeleteByID(w http.ResponseWriter, r *http.Request)
	BalanceSummary(w http.ResponseWriter, r *http.Request)
}

func NewHandler(
	service TransactionService,
	errHandler apiError.ErrorHandler,
) TransactionHandler {
	return &transactionHandler{
		service:    service,
		errHandler: errHandler,
	}
}

func (h *transactionHandler) BalanceSummary(w http.ResponseWriter, r *http.Request) {
	monthStr := handler.ReadStringParam(r, "month", "")
	yearStr := handler.ReadStringParam(r, "year", "")

	var start, end *time.Time
	if monthStr != "" && yearStr != "" {
		s, e, err := utils.ParseMonthYear(monthStr, yearStr)
		if err != nil {
			h.errHandler.BadRequestResponse(w, r, err)
			return
		}
		start = &s
		end = &e
	}

	summary, err := h.service.BalanceSummary(r.Context(), start, end)
	if err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	handler.Respond(w, r, http.StatusOK, summary, nil, h.errHandler)
}

func (h *transactionHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	categoryIDStr := handler.ReadStringParam(r, "category_id", "")
	startDateStr := handler.ReadStringParam(r, "start_date", "")
	endDateStr := handler.ReadStringParam(r, "end_date", "")
	search := handler.ReadStringParam(r, "search", "")

	f, err := handler.GetFilters(r, []string{"created_at", "-created_at", "amount", "-amount"})
	if err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	var categoryID *uuid.UUID
	if categoryIDStr != "" {
		id, err := uuid.Parse(categoryIDStr)
		if err != nil {
			h.errHandler.BadRequestResponse(w, r, fmt.Errorf("category_id inválido: %w", err))
			return
		}
		categoryID = &id
	}

	var startDate, endDate *time.Time
	if startDateStr != "" {
		t, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			h.errHandler.BadRequestResponse(w, r, fmt.Errorf("start_date inválido: %w", err))
			return
		}
		startDate = &t
	}
	if endDateStr != "" {
		t, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			h.errHandler.BadRequestResponse(w, r, fmt.Errorf("end_date inválido: %w", err))
			return
		}
		endDate = &t
	}

	var searchSlice []string
	if search != "" {
		searchSlice = []string{search}
	}

	models, metadata, err := h.service.FindAll(r.Context(), categoryID, startDate, endDate, f, searchSlice...)
	if err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	dtos := make([]*TransactionDTO, 0, len(models))
	for _, m := range models {
		dtos = append(dtos, m.ToDTO())
	}

	handler.Respond(
		w, r,
		http.StatusOK,
		utils.Envelope{
			"content":  dtos,
			"metadata": metadata,
		},
		nil,
		h.errHandler,
	)
}

func (h *transactionHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	id, ok := handler.ParseUUID(w, r, h.errHandler)
	if !ok {
		return
	}

	model, err := h.service.FindById(r.Context(), id)
	if err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	handler.Respond(w, r, http.StatusOK, model.ToDTO(), nil, h.errHandler)
}

func (h *transactionHandler) Save(w http.ResponseWriter, r *http.Request) {
	var dto TransactionDTO
	if err := handler.ReadJSON(w, r, &dto); err != nil {
		h.errHandler.BadRequestResponse(w, r, err)
		return
	}

	model, err := dto.ToModel()
	if err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	if err := h.service.Insert(r.Context(), nil, model); err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	handler.Respond(w, r, http.StatusCreated, model.ToDTO(), nil, h.errHandler)
}

func (h *transactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	var dto TransactionDTO
	if err := handler.ReadJSON(w, r, &dto); err != nil {
		h.errHandler.BadRequestResponse(w, r, err)
		return
	}

	model, err := dto.ToModel()
	if err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	if err := h.service.Update(r.Context(), nil, model); err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	handler.Respond(w, r, http.StatusOK, model.ToDTO(), nil, h.errHandler)
}

func (h *transactionHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id, ok := handler.ParseUUID(w, r, h.errHandler)
	if !ok {
		return
	}

	if err := h.service.DeleteById(r.Context(), nil, id); err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	handler.Respond(w, r, http.StatusNoContent, nil, nil, h.errHandler)
}
