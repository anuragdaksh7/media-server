package service

import (
	"context"
	"errors"
	"fileserver/internal/storage/provider"
	"io"
	"mime/multipart"
	"path/filepath"
)

type UploadService struct {
	Storage provider.StorageProvider

	MaxUploadSize int64
}

func NewUploadService(
	storage provider.StorageProvider,
	maxUploadSize int64,
) *UploadService {

	return &UploadService{
		Storage:       storage,
		MaxUploadSize: maxUploadSize,
	}
}

func (s *UploadService) Upload(
	ctx context.Context,
	basePath string,
	fileHeader *multipart.FileHeader,
) error {

	if fileHeader.Size > s.MaxUploadSize {
		return errors.New("file exceeds max upload size")
	}

	src, err := fileHeader.Open()

	if err != nil {
		return err
	}

	defer src.Close()

	dstPath := filepath.Join(
		basePath,
		fileHeader.Filename,
	)

	return s.Storage.Write(
		ctx,
		dstPath,
		src,
	)
}

func (s *UploadService) UploadStream(
	ctx context.Context,
	path string,
	reader io.Reader,
) error {

	return s.Storage.Write(
		ctx,
		path,
		reader,
	)
}
