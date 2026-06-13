package job

import (
	"log/slog"
	"net/http"
)

type JobHandler struct {
	jobService *JobService
}

func NewJobHandler(jobService *JobService) *JobHandler {
	return &JobHandler{jobService: jobService}
}

func (jh *JobHandler) GetJobStatus(w http.ResponseWriter, req *http.Request) {
	jobId := req.PathValue("id")
	slog.Info(jobId)
	jh.jobService.GetJobStatus(jobId)
}
