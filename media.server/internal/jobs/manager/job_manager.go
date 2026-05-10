package manager

import (
	model2 "fileserver/internal/auth/model"
	"fileserver/internal/jobs/model"
	"sync"
	"time"

	"github.com/google/uuid"
)

type JobManager struct {
	mu sync.RWMutex

	jobs map[string]*model.Job
}

func NewJobManager() *JobManager {

	return &JobManager{
		jobs: make(map[string]*model.Job),
	}
}

func (m *JobManager) CreateJob(
	jobType model.JobType,
	payload map[string]interface{},
) *model.Job {

	m.mu.Lock()
	defer m.mu.Unlock()

	job := &model.Job{
		Base: model2.Base{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},

		Type: jobType,

		Status: model.JobPending,

		Progress: 0,

		Payload: payload,
	}

	m.jobs[job.ID.String()] = job

	return job
}

func (m *JobManager) UpdateStatus(
	jobID string,
	status model.JobStatus,
) {

	m.mu.Lock()
	defer m.mu.Unlock()

	job, exists := m.jobs[jobID]

	if !exists {
		return
	}

	job.Status = status
	job.UpdatedAt = time.Now()
}

func (m *JobManager) UpdateProgress(
	jobID string,
	progress int,
) {

	m.mu.Lock()
	defer m.mu.Unlock()

	job, exists := m.jobs[jobID]

	if !exists {
		return
	}

	job.Progress = progress
	job.UpdatedAt = time.Now()
}

func (m *JobManager) FailJob(
	jobID string,
	err string,
) {

	m.mu.Lock()
	defer m.mu.Unlock()

	job, exists := m.jobs[jobID]

	if !exists {
		return
	}

	job.Status = model.JobFailed
	job.Error = err
	job.UpdatedAt = time.Now()
}

func (m *JobManager) GetJob(
	jobID string,
) (*model.Job, bool) {

	m.mu.RLock()
	defer m.mu.RUnlock()

	job, exists := m.jobs[jobID]

	return job, exists
}

func (m *JobManager) ListJobs() []*model.Job {

	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*model.Job, 0)

	for _, job := range m.jobs {
		result = append(result, job)
	}

	return result
}
