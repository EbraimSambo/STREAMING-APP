package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"stream/internal/tools"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func UploadVideo(c echo.Context) error {
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

	go tools.TranscodeToHLS(videoPath, videoFolder)

	return c.JSON(http.StatusOK, map[string]string{
		"message":  "Upload concluído com sucesso",
		"video_id": videoID,
		"playlist": fmt.Sprintf("/uploads/%s/master.m3u8", videoID),
	})
}
