package main

import (
	"github.com/joho/godotenv"

	"github.com/AkbarFikri/BREECE-BE/internal/app/config"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load env file")
	}

	db := config.NewDatabase()
	app := config.NewGin()

	config.StartUp(&config.StartUpConfig{
		DB:  db,
		App: app,
	})

	app.Run()
}
