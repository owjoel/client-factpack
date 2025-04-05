package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/clients/pkg/service"
)

type JobHandler struct {
	service *service.JobService
}

func NewJobHandler(service *service.JobService) *JobHandler {
	return &JobHandler{service: service}
}

func (h *JobHandler) GetJob(c *gin.Context) {
	jobID := c.Param("id")

	if jobID == "" {
		c.JSON(http.StatusBadRequest, model.StatusRes{Status: "Missing id"})
		return
	}

	job, err := h.service.GetJob(c.Request.Context(), jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, model.StatusRes{Status: fmt.Sprintf("Could not retrieve job: %v", err)})
		return
	}

	resp(c, http.StatusOK, job)
}

func (h *JobHandler) GetAllJobs(c *gin.Context) {
	query := &model.GetJobsQuery{}

	if err := c.ShouldBindQuery(query); err != nil {
		c.JSON(http.StatusBadRequest, model.StatusRes{Status: fmt.Sprintf("Invalid request: %v", err)})
		return
	}

	total, jobs, err := h.service.GetAllJobs(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.StatusRes{Status: fmt.Sprintf("Could not retrieve jobs: %v", err)})
		return
	}

	resp(c, http.StatusOK, model.GetJobsResponse{Total: total, Jobs: jobs})
}
