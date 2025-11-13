package main

import (
	"go-first-project/database"
	"go-first-project/handlers"
	helper "go-first-project/internal"
	"go-first-project/repository"
	"go-first-project/services"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	dbConfig := database.Config{
		Host:     helper.GetEnv("DB_HOST", ""),
		Port:     helper.GetEnv("DB_PORT", ""),
		User:     helper.GetEnv("DB_USER", ""),
		Password: helper.GetEnv("DB_PASSWORD", ""),
		DBName:   helper.GetEnv("DB_NAME", ""),
	}

	if err := database.Connect(dbConfig); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	if err := database.InitSchema(); err != nil {
		log.Fatal("Failed to initialize database schema:", err)
	}

	wRepo := repository.NewWeatherRepository(database.DB)
	wServices := services.NewWeatcherService(wRepo)
	wHandlers := handlers.WeatherHandler{WeatcherService: *wServices}

	http.HandleFunc("/weather", wHandlers.GetWeatherHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}

}
