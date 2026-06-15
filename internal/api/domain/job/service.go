package job

import (
	"log/slog"

	"github.com/hibiken/asynq"
)

type JobService struct {
	asynqInspector *asynq.Inspector
}

func NewJobService(asynqInspector *asynq.Inspector) *JobService {
	return &JobService{asynqInspector: asynqInspector}
}

func (js *JobService) GetJobStatus(jobId string) (*string, error) {
	info, err := js.asynqInspector.GetTaskInfo("ingest", jobId)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	status := info.State.String()
	return &status, nil
}
