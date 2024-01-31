package main

import (
	"log"
	"url-shortener-go/config"
	"url-shortener-go/routes"
	"url-shortener-go/scheduler"

	"github.com/labstack/echo/v4"
)

func main() {
	config.InitCache()
	config.InitDB()

	go scheduler.Init()

	e := echo.New()
	routes.Init(e)

	if err := e.Start(":8000"); err != nil {
		log.Fatal(err)
	}
}
