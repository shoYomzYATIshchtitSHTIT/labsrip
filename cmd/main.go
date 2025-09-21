package main

import (
	"log"

	"Backend-RIP/internal/api"
)

func main() {
	log.Println("Application start")
	api.StartServer()
	log.Println("Application terminated")
}
