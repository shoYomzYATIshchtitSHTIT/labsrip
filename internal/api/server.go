package api

import (
	"Backend-RIP/internal/app/handler"
	"Backend-RIP/internal/app/repository"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("Ошибка инициализация репозитория")
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	r.GET("/", handler.GetIntervals)
	r.GET("/interval/:id", handler.GetInterval)
	r.GET("/cart", handler.GetCart)

	r.Run()
	log.Println("Server down")
}
