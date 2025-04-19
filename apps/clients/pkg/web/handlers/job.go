package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	errorx "github.com/owjoel/client-factpack/apps/clients/pkg/api/errors"
	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/clients/pkg/service"
)

type JobHandler struct {
	service service.JobServiceInterface
}

func NewJobHandler(service service.JobServiceInterface) *JobHandler {
	return &JobHandler{service: service}
}

func (h *JobHandler) GetJob(c *gin.Context) {
	jobID := c.Param("id")

	if jobID == "" {
		resp(c, http.StatusBadRequest, model.ErrorResponse{Message: "Missing id"})
		return
	}

	job, err := h.service.GetJob(c.Request.Context(), jobID)
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidInput):
			resp(c, http.StatusBadRequest, model.ErrorResponse{Message: "Invalid job ID"})
		case errors.Is(err, errorx.ErrNotFound):
			resp(c, http.StatusNotFound, model.ErrorResponse{Message: "Job not found"})
		default:
			resp(c, http.StatusInternalServerError, model.ErrorResponse{Message: "Internal server error"})
		}
		return
	}

	resp(c, http.StatusOK, job)
}

func (h *JobHandler) GetAllJobs(c *gin.Context) {
	query := &model.GetJobsQuery{}

	if err := c.ShouldBindQuery(query); err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidInput):
			resp(c, http.StatusBadRequest, model.ErrorResponse{Message: "Invalid request"})
		default:
			resp(c, http.StatusInternalServerError, model.ErrorResponse{Message: "Internal server error"})
		}
		return
	}

	total, jobs, err := h.service.GetAllJobs(c.Request.Context(), query)
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidInput):
			resp(c, http.StatusBadRequest, model.ErrorResponse{Message: "Invalid request"})
		case errors.Is(err, errorx.ErrDependencyFailed):
			resp(c, http.StatusServiceUnavailable, model.ErrorResponse{Message: "Database error — please try again later"})
		default:
			resp(c, http.StatusInternalServerError, model.ErrorResponse{Message: "Internal server error"})
		}
		return
	}

	resp(c, http.StatusOK, model.GetJobsResponse{Total: total, Jobs: jobs})
}
