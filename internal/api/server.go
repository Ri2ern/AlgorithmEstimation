package api

import (
	"lab1-algorithms/internal/app/handler"
	"lab1-algorithms/internal/app/repository"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	log.Println("Application start!")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("Ошибка инициализации репозитория")
	}

	h := handler.NewHandler(repo)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	r.GET("/feed", h.FeedHandler)
	r.GET("/add", h.AddHandler)
	r.GET("/grid", h.GridHandler)
	r.GET("/", func(c *gin.Context) {
		c.Redirect(303, "/feed")
	})

	log.Println("Server started on http://localhost:8080")
	r.Run(":8080")

	log.Println("Application terminated!")
}