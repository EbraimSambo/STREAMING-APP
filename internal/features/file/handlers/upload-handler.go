package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"stream/ent"
	"stream/internal/features/file/repository"
	"stream/internal/features/file/service"
	"stream/internal/tools"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func UploadVideo(c echo.Context, client *ent.Client) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Arquivo é obrigatório"})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao abrir arquivo"})
	}
	defer src.Close()

	videoID := uuid.New().String()
	videoFolder := filepath.Join("uploads", videoID)
	err = os.MkdirAll(videoFolder, 0755)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao criar pasta de destino"})
	}

	videoPath := filepath.Join(videoFolder, "original.mp4")
	dst, err := os.Create(videoPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao salvar arquivo"})
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao copiar arquivo"})
	}

	repo := repository.NewFileRepository(client)
	service := service.NewFileService(repo)
	_, err = service.SaveFile(context.Background(), videoID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Erro ao guardar video",
			"msg":   err.Error(),
		})
	}
	go tools.TranscodeToHLS(videoPath, videoFolder, videoID, client)

	return c.JSON(http.StatusOK, map[string]string{
		"message":  "Upload concluído com sucesso",
		"video_id": videoID,
		"playlist": fmt.Sprintf("/uploads/%s/master.m3u8", videoID),
	})
}
