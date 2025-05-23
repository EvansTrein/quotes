package controller

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	customErr "quotes/pkg/error"
)

// Universal structure for sending responses
type BaseControllerResponce struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Status  int    `json:"status"`
}

type BaseController struct {
	Log *slog.Logger
}

type BaseControllerDeps struct {
	*slog.Logger
}

func NewBaseController(deps *BaseControllerDeps) *BaseController {
	return &BaseController{Log: deps.Logger}
}

func (h *BaseController) SendJsonResp(w http.ResponseWriter, status int, data any) {
	jsonResponse, err := json.Marshal(data)
	if err != nil {
		h.Log.Error("failed to marshal JSON", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(jsonResponse); err != nil {
		h.Log.Error("!!ATTENTION!! failed to write JSON response", "error", err)
	}
}

func (h *BaseController) HandleError(w http.ResponseWriter, err error) {
	op := "BaseController: HandleError func"
	log := h.Log.With(slog.String("operation", op))

	switch {
	case errors.Is(err, customErr.ErrNoFields):
		log.Warn("validation failed", "error", err)
		h.SendJsonResp(w, 400, &BaseControllerResponce{
			Status:  http.StatusBadRequest,
			Message: "validation failed",
			Error:   err.Error(),
		})
	case errors.Is(err, customErr.ErrInvalidTypeID):
		log.Warn("incorrect id", "error", err)
		h.SendJsonResp(w, 400, &BaseControllerResponce{
			Status:  http.StatusBadRequest,
			Message: "incorrect id",
			Error:   err.Error(),
		})
	case errors.Is(err, customErr.ErrRecordNotFound):
		log.Warn("fail find the record", "error", err)
		h.SendJsonResp(w, 404, &BaseControllerResponce{
			Status:  http.StatusNotFound,
			Message: "fail find the record",
			Error:   err.Error(),
		})
	case errors.Is(err, customErr.ErrNoQuotesAvailable):
		log.Warn("database empty", "error", err)
		h.SendJsonResp(w, 404, &BaseControllerResponce{
			Status:  http.StatusNotFound,
			Message: "database empty",
			Error:   err.Error(),
		})
	case errors.Is(err, customErr.ErrAuthorNotFound):
		log.Warn("fail find the author", "error", err)
		h.SendJsonResp(w, 404, &BaseControllerResponce{
			Status:  http.StatusNotFound,
			Message: "fail find the author",
			Error:   err.Error(),
		})
	case errors.Is(err, context.DeadlineExceeded):
		log.Error("request processing exceeded the allowed time limit", "error", err)
		h.SendJsonResp(w, 504, &BaseControllerResponce{
			Status:  http.StatusGatewayTimeout,
			Message: "request processing exceeded the allowed time limit",
			Error:   err.Error(),
		})
	default:
		log.Error("internal server error", "error", err)
		h.SendJsonResp(w, 500, &BaseControllerResponce{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
			Error:   err.Error(),
		})
	}
}
