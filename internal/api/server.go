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

	// Основной репозиторий интервалов
	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("Ошибка инициализации репозитория")
	}

	// Репозиторий заявок на композиции
	compositionRepo, err := repository.NewCompositionRequestRepository()
	if err != nil {
		logrus.Error("Ошибка инициализации репозитория заявок")
	}

	// Создаём обработчик с обоими репозиториями
	h := handler.NewHandler(repo, compositionRepo)

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	r.GET("/", h.GetIntervals)
	r.GET("/interval/:id", h.GetInterval)
	r.GET("/composition/:id", h.GetComposition) // теперь ID через путь

	r.Run()
	log.Println("Server down")
}
