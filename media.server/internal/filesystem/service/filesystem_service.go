package service

import (
	"context"
	"fileserver/internal/storage/provider"
)

type FilesystemService struct {
	Storage provider.StorageProvider
}

func NewFilesystemService(
	storage provider.StorageProvider,
) *FilesystemService {

	return &FilesystemService{
		Storage: storage,
	}
}

func (s *FilesystemService) List(
	ctx context.Context,
	path string,
) ([]provider.FileInfo, error) {

	return s.Storage.List(ctx, path)
}

func (s *FilesystemService) Stat(
	ctx context.Context,
	path string,
) (*provider.FileInfo, error) {

	return s.Storage.Stat(ctx, path)
}

func (s *FilesystemService) Delete(
	ctx context.Context,
	path string,
) error {

	return s.Storage.Delete(ctx, path)
}

func (s *FilesystemService) Move(
	ctx context.Context,
	src string,
	dst string,
) error {

	return s.Storage.Move(ctx, src, dst)
}

func (s *FilesystemService) Copy(
	ctx context.Context,
	src string,
	dst string,
) error {

	return s.Storage.Copy(ctx, src, dst)
}

func (s *FilesystemService) Mkdir(
	ctx context.Context,
	path string,
) error {

	return s.Storage.Mkdir(ctx, path)
}
