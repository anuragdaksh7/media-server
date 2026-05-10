package service

import (
	"fileserver/internal/jobs/manager"
	"fileserver/internal/jobs/model"
)

type JobService struct {
	Manager *manager.JobManager
}

func NewJobService(
	manager *manager.JobManager,
) *JobService {

	return &JobService{
		Manager: manager,
	}
}

func (s *JobService) GetJob(
	jobID string,
) (*model.Job, bool) {

	return s.Manager.GetJob(jobID)
}

func (s *JobService) ListJobs() []*model.Job {

	return s.Manager.ListJobs()
}
