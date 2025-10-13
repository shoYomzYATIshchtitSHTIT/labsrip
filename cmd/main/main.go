package main

import (
	"Backend-RIP/internal/app/config"
	"Backend-RIP/internal/app/repository"
	"Backend-RIP/internal/pkg"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	router := gin.Default()
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	// Инициализируем репозиторий
	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Fatalf("error initializing repository: %v", err)
	}

	// Создаем приложение
	application := pkg.NewApp(conf, router, repo)

	// Запускаем приложение
	application.RunApp()
}
