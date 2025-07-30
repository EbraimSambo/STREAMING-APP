package service

import (
	"context"
	"stream/ent"
	core "stream/internal/core/pagination"
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
func (s *FileService) ChangeVisibility(ctx context.Context, fileRef string) (*string, error) {
	return s.Repo.ChangeVisibility(ctx, fileRef)
}

func (s *FileService) GetFolderFiles(d repository.DataParamsFiles) (*core.Pagination[ent.File], error) {
	return s.Repo.GetFolderFiles(d)
}

func (s *FileService) GetFile(ctx context.Context, fileRef string) (*ent.File, error) {
	return s.Repo.GetFile(ctx, fileRef)
}
