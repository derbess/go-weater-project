package main

import (
	"go-first-project/client"
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
		log.Printf("No .env file found, using system environment variables")
	}

	dbConfig := database.Config{
		Host:     helper.GetEnv("DB_HOST", ""),
		Port:     helper.GetEnv("DB_PORT", ""),
		User:     helper.GetEnv("DB_USER", ""),
		Password: helper.GetEnv("DB_PASSWORD", ""),
		DBName:   helper.GetEnv("DB_NAME", ""),
	}

	pool, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close(pool)

	if err := database.InitSchema(pool); err != nil {
		log.Fatal("Failed to initialize database schema:", err)
	}

	weatherClient := client.NewClient(
		helper.GetEnv("API_KEY_OPEN_WEATHER", ""),
		helper.GetEnv("API_KEY_WEATHER_API", ""),
	)

	wRepo := repository.NewWeatherRepository(pool)
	wServices := services.NewWeatcherService(wRepo, weatherClient)
	wHandlers := handlers.WeatherHandler{WeatcherService: *wServices}

	http.HandleFunc("/weather", wHandlers.GetWeatherHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Failed to start server:", err)
	}

}
