package category

import (
	"controle_financas/internal/core/handler"
	"controle_financas/pkg/httpjson"
	"controle_financas/pkg/httputil"
	"net/http"
)

type errorHandler interface {
	HandlerError(w http.ResponseWriter, r *http.Request, err error)
	ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error)
	BadRequestResponse(w http.ResponseWriter, r *http.Request, err error)
}

type CategoyHandler struct {
	service    categoryService
	errHandler errorHandler
}

type categoyHandler interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	FindByID(w http.ResponseWriter, r *http.Request)
	Save(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	DeleteByID(w http.ResponseWriter, r *http.Request)
}

func NewHandler(
	service categoryService,
	errHandler errorHandler,
) *CategoyHandler {
	return &CategoyHandler{
		service:    service,
		errHandler: errHandler,
	}
}

func (h *CategoyHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	name := httputil.ReadStringParam(r, "name", "")
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
		httpjson.Envelope{
			"content":  dtos,
			"metadata": metadata,
		},
		nil,
		h.errHandler,
	)
}

func (h *CategoyHandler) FindByID(w http.ResponseWriter, r *http.Request) {
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

func (h *CategoyHandler) Save(w http.ResponseWriter, r *http.Request) {
	var dto CategoryDTO
	if err := httputil.ReadJSON(w, r, &dto); err != nil {
		h.errHandler.BadRequestResponse(w, r, err)
		return
	}

	model, err := dto.ToModel()
	if err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	if err := h.service.Insert(r.Context(), model); err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	handler.Respond(w, r, http.StatusCreated, model.ToDTO(), nil, h.errHandler)
}

func (h *CategoyHandler) Update(w http.ResponseWriter, r *http.Request) {
	var dto CategoryDTO
	if err := httputil.ReadJSON(w, r, &dto); err != nil {
		h.errHandler.BadRequestResponse(w, r, err)
		return
	}

	model, err := dto.ToModel()
	if err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	if err := h.service.Update(r.Context(), model); err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	handler.Respond(w, r, http.StatusOK, model.ToDTO(), nil, h.errHandler)
}

func (h *CategoyHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id, ok := handler.ParseUUID(w, r, h.errHandler)
	if !ok {
		return
	}

	if err := h.service.DeleteById(r.Context(), id); err != nil {
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
