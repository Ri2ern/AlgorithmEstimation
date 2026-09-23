package api

import (
	"lab1-algorithms/internal/app/handler"
	"lab1-algorithms/internal/app/repository"
	"log"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	log.Println("Application start!")

	repo, err := repository.NewRepository()
	if err != nil {
		log.Fatal("Ошибка инициализации БД:", err)
	}

	h := handler.NewHandler(repo)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	// 3 GET метода
	r.GET("/feed", h.GetFeed)
	r.GET("/grid", h.GetGrid)
	r.GET("/add", h.GetAdd)
	
	// 3 POST метода
	r.POST("/add", h.PostCreateDraft)
	r.POST("/publish", h.PostPublish)
	r.POST("/delete", h.PostDelete)
	
	r.GET("/", func(c *gin.Context) {
		c.Redirect(303, "/feed")
	})

	log.Println("Server started on http://localhost:8080")
	r.Run(":8080")
}