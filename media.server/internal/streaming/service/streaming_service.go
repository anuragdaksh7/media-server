package service

import (
	"context"
	"fileserver/internal/storage/provider"
	"io"
)

type StreamingService struct {
	Storage provider.StorageProvider
}

func NewStreamingService(
	storage provider.StorageProvider,
) *StreamingService {

	return &StreamingService{
		Storage: storage,
	}
}

func (s *StreamingService) OpenFile(
	ctx context.Context,
	path string,
) (io.ReadSeekCloser, *provider.FileInfo, error) {

	file, err := s.Storage.Read(ctx, path)

	if err != nil {
		return nil, nil, err
	}

	info, err := s.Storage.Stat(ctx, path)

	if err != nil {
		file.Close()
		return nil, nil, err
	}

	return file, info, nil
}
