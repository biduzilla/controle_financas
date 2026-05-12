package category

import (
	"controle_financas/internal/core/domain/apiError"
	"controle_financas/internal/core/handler"
	"controle_financas/utils"
	"net/http"
)

type categoyHandler struct {
	service    CategoryService
	errHandler apiError.ErrorHandler
}

type CategoyHandler interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	FindByID(w http.ResponseWriter, r *http.Request)
	Save(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	DeleteByID(w http.ResponseWriter, r *http.Request)
}

func NewHandler(
	service CategoryService,
	errHandler apiError.ErrorHandler,
) CategoyHandler {
	return &categoyHandler{
		service:    service,
		errHandler: errHandler,
	}
}

func (h *categoyHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	name := handler.ReadStringParam(r, "name", "")
	f, err := handler.GetFilters(r, []string{"nome", "-nome"})
	if err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	models, metadata, err := h.service.FindAll(r.Context(), f, name)

	if err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	dtos := make([]*CategoryDTO, 0, len(models))
	for _, m := range models {
		dtos = append(dtos, m.ToDTO())
	}

	handler.Respond(
		w,
		r,
		http.StatusOK,
		utils.Envelope{
			"content":  dtos,
			"metadata": metadata,
		},
		nil,
		h.errHandler,
	)
}

func (h *categoyHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	id, ok := handler.ParseUUID(w, r, h.errHandler)
	if !ok {
		return
	}

	model, err := h.service.FindById(r.Context(), id)
	if err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	handler.Respond(
		w, r,
		http.StatusOK,
		model.ToDTO(),
		nil,
		h.errHandler,
	)
}

func (h *categoyHandler) Save(w http.ResponseWriter, r *http.Request) {
	var dto CategoryDTO
	if err := handler.ReadJSON(w, r, &dto); err != nil {
		h.errHandler.BadRequestResponse(w, r, err)
		return
	}

	model, err := dto.ToModel()
	if err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	if err := h.service.Insert(r.Context(), model, nil); err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	handler.Respond(w, r, http.StatusCreated, model.ToDTO(), nil, h.errHandler)
}

func (h *categoyHandler) Update(w http.ResponseWriter, r *http.Request) {
	var dto CategoryDTO
	if err := handler.ReadJSON(w, r, &dto); err != nil {
		h.errHandler.BadRequestResponse(w, r, err)
		return
	}

	model, err := dto.ToModel()
	if err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	if err := h.service.Update(r.Context(), model, nil); err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	handler.Respond(w, r, http.StatusOK, model.ToDTO(), nil, h.errHandler)
}

func (h *categoyHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id, ok := handler.ParseUUID(w, r, h.errHandler)
	if !ok {
		return
	}

	if err := h.service.DeleteById(r.Context(), id, nil); err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	handler.Respond(
		w,
		r,
		http.StatusNoContent,
		nil,
		nil,
		h.errHandler,
	)
}
