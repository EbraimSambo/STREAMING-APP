package repository

import (
	"context"
	"log"
	"stream/ent"
	"stream/ent/file"
	core "stream/internal/core/pagination"
)

type FileRepository struct {
	Client *ent.Client
}

type DataUploadFIle struct {
	File string `json:"file"`
}

func NewFileRepository(clint *ent.Client) FileRepository {
	return FileRepository{Client: clint}
}

func (r FileRepository) SaveFile(ctx context.Context, videoID, fileName string) (*string, error) {

	tx, err := r.Client.Tx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	newFile, err := tx.File.
		Create().
		SetID(videoID).
		SetFileName(fileName).
		SetStatus("PENDING").
		Save(ctx)

	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	return &newFile.ID, nil
}

func (r FileRepository) ChangeVisibility(ctx context.Context, videoID string) (*string, error) {
	tx, err := r.Client.Tx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	_, err = tx.File.
		Update().
		Where(
			file.IDEQ(videoID),
		).
		SetVisibility(true).
		Save(ctx)

	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	return &videoID, nil

}

type DataParamsFiles struct {
	Ctx            context.Context
	Client         *ent.Client
	DataPagination core.DataPagination
}

func (r *FileRepository) GetFolderFiles(d DataParamsFiles) (*core.Pagination[ent.File], error) {
	offset := (d.DataPagination.Page - 1) * d.DataPagination.Limit

	total, err := d.Client.File.
		Query().
		Where(
			file.VisibilityEQ(true),
		).
		Count(d.Ctx)

	if err != nil {
		return nil, err
	}

	files, err := d.Client.File.
		Query().
		Limit(d.DataPagination.Limit).
		Offset(offset).
		Where(
			file.VisibilityEQ(true),
		).
		Order(ent.Desc(file.FieldCreatedAt)).All(d.Ctx)

	if err != nil {
		return nil, err
	}

	totalItems := total
	totalPages := (totalItems + d.DataPagination.Limit - 1) / d.DataPagination.Limit

	var prevPage *int
	var nextPage *int

	if d.DataPagination.Page > 1 {
		p := d.DataPagination.Page - 1
		prevPage = &p
	}
	if d.DataPagination.Page < totalPages {
		n := d.DataPagination.Page + 1
		nextPage = &n
	}

	return &core.Pagination[ent.File]{
		Items: files,
		Info: core.PageInfo{
			Total:       totalItems,
			Page:        d.DataPagination.Page,
			PerPage:     d.DataPagination.Limit,
			TotalPages:  totalPages,
			PrevPage:    prevPage,
			NextPage:    nextPage,
			HasNextPage: nextPage != nil,
		},
	}, nil
}

func (r *FileRepository) GetFile(ctx context.Context, fileRef string) (*ent.File, error) {
	return r.Client.File.Query().Where(
		file.IDEQ(fileRef),
		file.VisibilityEQ(true),
	).First(ctx)
}