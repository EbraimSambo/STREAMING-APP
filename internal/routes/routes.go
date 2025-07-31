package routes

import (
	"stream/ent"
	"stream/internal/features/file/handlers"

	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Echo, client *ent.Client) {
	e.POST("/upload", func(c echo.Context) error {
		return handlers.UploadVideo(c, client)
	})

	e.GET("/videos", func(c echo.Context) error {
		return handlers.GetVideos(c, client)
	})

	e.GET("/videos/single/:fileRef", func(c echo.Context) error {
		return handlers.GetVideo(c, client)
	})

	// Serve os vídeos via rota /videos/<uuid>/<arquivo>
	e.Static("/videos", "uploads")
}
