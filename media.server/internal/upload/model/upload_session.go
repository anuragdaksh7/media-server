package model

import (
	"fileserver/internal/auth/model"
	"time"
)

type UploadStatus string

const (
	UploadPending   UploadStatus = "pending"
	UploadRunning   UploadStatus = "running"
	UploadCompleted UploadStatus = "completed"
	UploadFailed    UploadStatus = "failed"
)

type UploadSession struct {
	model.Base
	FileName      string       `json:"file_name"`
	Path          string       `json:"path"`
	Size          int64        `json:"size"`
	UploadedBytes int64        `json:"uploaded_bytes"`
	Status        UploadStatus `json:"status"`
	CreatedAt     time.Time    `json:"created_at"`
}
