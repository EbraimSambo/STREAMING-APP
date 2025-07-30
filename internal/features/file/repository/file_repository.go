package repository

import (
	"context"
	"log"
	"stream/ent"
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

func (r FileRepository) SaveFile(ctx context.Context, file string) (*string, error) {

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
		SetFile(file).
		Save(ctx)

	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	return &newFile.File, nil
}
