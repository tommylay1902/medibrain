package job

import (
	"encoding/json"
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
	status, err := jh.jobService.GetJobStatus(jobId)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"status": *status,
	}); err != nil {
		slog.Error("Error writing jobId json result", slog.Any("err", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
