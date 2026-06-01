package handler

import (
	"controle_financas/internal/core/domain/apiError"
	"controle_financas/internal/core/filters"
	"controle_financas/internal/core/validator"
	"controle_financas/pkg/httpjson"
	"controle_financas/pkg/httputil"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func ParseIntID(
	w http.ResponseWriter,
	r *http.Request,
	errRsp apiError.ErrorHandler,
) (int64, bool) {
	id, err := readIntPathVariable(r, "id")
	if err != nil {
		errRsp.BadRequestResponse(w, r, err)
		return 0, false
	}
	return id, true
}

func GetFilters(r *http.Request, sortSafelist []string) (filters.Filters, error) {
	v := validator.New()

	var f = filters.Filters{}
	f.Page = httputil.ReadIntParam(r, "page", 1, v)
	f.PageSize = httputil.ReadIntParam(r, "page_size", 20, v)
	f.Sort = httputil.ReadStringParam(r, "sort", "")
	f.SortSafelist = sortSafelist

	if filters.ValidateFilters(v, f); !v.Valid() {
		return filters.Filters{}, apiError.NewValidationError(v.Errors)
	}

	return f, nil

}

func ParseStringField(
	w http.ResponseWriter,
	r *http.Request,
	errRsp apiError.ErrorHandler,
	field string,
) (string, bool) {
	value, err := readStringPathVariable(r, field)
	if err != nil {
		errRsp.BadRequestResponse(w, r, err)
		return "", false
	}
	return value, true
}

func ParseUUID(
	w http.ResponseWriter,
	r *http.Request,
	errRsp apiError.ErrorHandler,
) (uuid.UUID, bool) {

	id, err := readStringPathVariable(r, "id")
	if err != nil {
		errRsp.BadRequestResponse(w, r, err)
		return uuid.Nil, false
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		errRsp.BadRequestResponse(w, r, err)
		return uuid.Nil, false
	}

	return uid, true
}

func Respond(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	data any,
	headers http.Header,
	errRsp apiError.ErrorHandler,
) {
	err := httpjson.WriteJSON(w, status, data, headers)
	if err != nil {
		errRsp.ServerErrorResponse(w, r, err)
		return
	}
}

func readIntPathVariable(r *http.Request, key string) (int64, error) {
	s := chi.URLParam(r, key)

	if s == "" {
		return 0, fmt.Errorf("missing path parameter: %s", key)
	}

	value, err := strconv.ParseInt(s, 10, 64)

	if err != nil {
		return 0, fmt.Errorf("invalid %s parameter", key)
	}

	return value, nil
}

func readStringPathVariable(r *http.Request, key string) (string, error) {
	s := chi.URLParam(r, key)

	if s == "" {
		return "", fmt.Errorf("missing path parameter: %s", key)
	}

	return s, nil
}
