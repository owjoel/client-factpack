package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	errorx "github.com/owjoel/client-factpack/apps/clients/pkg/api/errors"
	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/clients/pkg/service"
)

type LogHandler struct {
	logService service.LogServiceInterface
}

func NewLogHandler(logService service.LogServiceInterface) *LogHandler {
	return &LogHandler{logService: logService}
}

func (h *LogHandler) GetLogs(c *gin.Context) {
	query := &model.GetLogsQuery{}
	if err := c.ShouldBindQuery(query); err != nil {
		resp(c, http.StatusBadRequest, model.ErrorResponse{Message: "Invalid query parameters"})
		return
	}

	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize <= 0 || query.PageSize > 100 {
		query.PageSize = 20
	}

	total, logs, err := h.logService.GetLogs(c.Request.Context(), query)
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidInput):
			resp(c, http.StatusBadRequest, model.ErrorResponse{Message: "Invalid filters or client ID"})
		case errors.Is(err, errorx.ErrDependencyFailed):
			resp(c, http.StatusServiceUnavailable, model.ErrorResponse{Message: "Database error — please try again later"})
		default:
			resp(c, http.StatusInternalServerError, model.ErrorResponse{Message: "Internal server error"})
		}
		return
	}

	resp(c, http.StatusOK, model.GetLogsResponse{
		Total: total,
		Logs:  logs,
	})
}

func (h *LogHandler) GetLog(c *gin.Context) {
	logID := c.Param("id")
	if logID == "" {
		resp(c, http.StatusBadRequest, model.ErrorResponse{Message: "Log ID is required"})
		return
	}

	newLog, err := h.logService.GetLog(c.Request.Context(), logID)
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidInput):
			resp(c, http.StatusBadRequest, model.ErrorResponse{Message: "Invalid log ID format"})
		case errors.Is(err, errorx.ErrNotFound):
			resp(c, http.StatusNotFound, model.ErrorResponse{Message: "Log not found"})
		default:
			resp(c, http.StatusInternalServerError, model.ErrorResponse{Message: "Internal server error"})
		}
		return
	}

	resp(c, http.StatusOK, model.GetLogResponse{
		Log: newLog,
	})
}


func (h *LogHandler) CreateLog(c *gin.Context) {
	newLog := &model.Log{}
	if err := c.ShouldBindJSON(newLog); err != nil {
		resp(c, http.StatusBadRequest, model.ErrorResponse{Message: "Invalid log data"})
		return
	}

	id, err := h.logService.CreateLog(c.Request.Context(), newLog)
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidInput):
			resp(c, http.StatusBadRequest, model.ErrorResponse{Message: "Invalid log data"})
		case errors.Is(err, errorx.ErrDependencyFailed):
			resp(c, http.StatusServiceUnavailable, model.ErrorResponse{Message: "Database error — please try again later"})
		default:
			resp(c, http.StatusInternalServerError, model.ErrorResponse{Message: "Internal server error"})
		}
		return
	}

	resp(c, http.StatusCreated, model.CreateLogResponse{ID: id})
}

