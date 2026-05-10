package provider

import (
	"context"
	"io"
	"time"
)

type FileInfo struct {
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	Size       int64     `json:"size"`
	MimeType   string    `json:"mime_type"`
	IsDir      bool      `json:"is_dir"`
	ModifiedAt time.Time `json:"modified_at"`
}

type StorageProvider interface {
	List(
		ctx context.Context,
		path string,
	) ([]FileInfo, error)

	Stat(
		ctx context.Context,
		path string,
	) (*FileInfo, error)

	Read(
		ctx context.Context,
		path string,
	) (io.ReadSeekCloser, error)

	Write(
		ctx context.Context,
		path string,
		reader io.Reader,
	) error

	Delete(
		ctx context.Context,
		path string,
	) error

	Move(
		ctx context.Context,
		src string,
		dst string,
	) error

	Copy(
		ctx context.Context,
		src string,
		dst string,
	) error

	Mkdir(
		ctx context.Context,
		path string,
	) error
}
