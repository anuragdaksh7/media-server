package model

import (
	"fileserver/internal/auth/model"
)

type JobStatus string

const (
	JobPending   JobStatus = "pending"
	JobRunning   JobStatus = "running"
	JobCompleted JobStatus = "completed"
	JobFailed    JobStatus = "failed"
)

type JobType string

const (
	JobMove   JobType = "move"
	JobCopy   JobType = "copy"
	JobDelete JobType = "delete"
	JobUpload JobType = "upload"
)

type Job struct {
	model.Base
	Type     JobType                `json:"type"`
	Status   JobStatus              `json:"status"`
	Progress int                    `json:"progress"`
	Error    string                 `json:"error,omitempty"`
	Payload  map[string]interface{} `json:"payload"`
}
