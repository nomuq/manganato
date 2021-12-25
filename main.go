package main

import (
	"net/http"
	"os"

	nato "github.com/kaikaew13/manganato-api"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlash())

	searcher := nato.NewSearcher()

	// Routes
	e.GET("/", func(c echo.Context) error {
		mangas, err := searcher.SearchLatestUpdatedManga()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, mangas)
	})

	e.GET("/manga/:id", func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return c.JSON(http.StatusBadRequest, "id is required")
		}
		manga, err := searcher.PickManga(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, manga)
	})

	e.GET("/manga/:id/chapter/:ch", func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return c.JSON(http.StatusBadRequest, "id is required")
		}
		ch := c.Param("ch")
		if ch == "" {
			return c.JSON(http.StatusBadRequest, "chapter is required")
		}

		chapter, err := searcher.ReadMangaChapter(id, ch)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, chapter)
	})

	// Start server
	e.Logger.Fatal(e.Start(":" + port))
}
