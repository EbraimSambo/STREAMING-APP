package service

import (
	"context"
	"stream/internal/features/file/repository"
)

type FileService struct {
	Repo repository.FileRepository
}

func NewFileService(repo repository.FileRepository) FileService {
	return FileService{Repo: repo}
}

func (s *FileService) SaveFile(ctx context.Context, file string) (*string, error) {
	return s.Repo.SaveFile(ctx, file)
}
