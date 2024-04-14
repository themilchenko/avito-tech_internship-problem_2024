package main

import (
	"flag"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/themilchenko/avito_internship-problem_2024/internal/app"
	"github.com/themilchenko/avito_internship-problem_2024/internal/config"
)

func main() {
	var configPath string
	config.ParseFlag(&configPath)
	flag.Parse()

	cfg := config.New()
	if err := cfg.Open(configPath); err != nil {
		log.Fatal("Failed to open config file")
	}

	e := echo.New()
	app := app.NewServer(e, cfg)
	if err := app.Start(); err != nil {
		app.GetEchoInstance().Logger.Error(err)
	}
}
