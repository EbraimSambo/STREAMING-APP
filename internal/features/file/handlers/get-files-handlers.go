package handlers

import (
	"net/http"
	"strconv"
	"stream/ent"
	core "stream/internal/core/pagination"
	"stream/internal/features/file/repository"
	"stream/internal/features/file/service"

	"github.com/labstack/echo/v4"
)

func GetVideos(c echo.Context, client *ent.Client) error {
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")

	page, err := strconv.Atoi(pageParam)
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 5
	}

	repo := repository.NewFileRepository(client)
	service := service.NewFileService(repo)

	result, err := service.GetFolderFiles(repository.DataParamsFiles{
		Ctx:    c.Request().Context(),
		Client: client,
		DataPagination: core.DataPagination{
			Page:  page,
			Limit: limit,
		},
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao buscar dados"})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": result,
	})
}

func GetVideo(c echo.Context, client *ent.Client) error {
	fileRef := c.Param("fileRef")
	repo := repository.NewFileRepository(client)
	service := service.NewFileService(repo)
	result, err := service.GetFile(c.Request().Context(), fileRef)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao buscar dados"})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": result,
	})

}
