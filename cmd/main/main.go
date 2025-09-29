package main

import (
	"fmt"
	"html/template"

	"Backend-RIP/internal/app/config"
	"Backend-RIP/internal/app/dsn"
	"Backend-RIP/internal/app/handler"
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

	router.SetFuncMap(template.FuncMap{
		"find_amperage": func(a, b float64) float64 {
			return a * 1000 / b
		},
	})

	postgresString := dsn.FromEnv()
	fmt.Println(postgresString)

	rep, errRep := repository.New(postgresString)
	if errRep != nil {
		logrus.Fatalf("error initializing repository: %v", errRep)
	}

	hand := handler.NewHandler(rep)

	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}
