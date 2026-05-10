package service

import (
	"context"
	"fileserver/internal/jobs/manager"
	"fileserver/internal/jobs/model"
)

func (s *FilesystemService) MoveAsync(
	ctx context.Context,
	jobManager *manager.JobManager,
	src string,
	dst string,
) *model.Job {

	job := jobManager.CreateJob(
		model.JobMove,
		map[string]interface{}{
			"source":      src,
			"destination": dst,
		},
	)

	go func() {
		jobID := job.ID.String()

		jobManager.UpdateStatus(
			jobID,
			model.JobRunning,
		)

		err := s.Storage.Move(
			ctx,
			src,
			dst,
		)

		if err != nil {

			jobManager.FailJob(
				jobID,
				err.Error(),
			)

			return
		}

		jobManager.UpdateProgress(
			jobID,
			100,
		)

		jobManager.UpdateStatus(
			jobID,
			model.JobCompleted,
		)
	}()

	return job
}
