package services

import (
	"context"
	"fmt"
	"go-first-project/internal"
	"go-first-project/models"
	"go-first-project/repository"
	"net/http"
	"time"
)

type WeatcherService struct {
	repo *repository.WeatherRepository
}

func NewWeatcherService(repo *repository.WeatherRepository) *WeatcherService {
	return &WeatcherService{repo: repo}
}

func (s *WeatcherService) GetWeather(ctx context.Context, client *http.Client, city string) (*models.WeatherApiResponse, error) {
	start := time.Now()

	openWeatherChan := make(chan *models.OpenWeather)
	weatherApiChan := make(chan *models.WeatherApi)

	go func() {
		w, _ := internal.FeatchOpenWeather(ctx, client, city)
		fmt.Println("OpenWeather time:", time.Since(start))
		openWeatherChan <- w
	}()

	go func() {
		w2, _ := internal.FetchWeatherApi(ctx, client, city)
		fmt.Println("WeatherApi time:", time.Since(start))
		weatherApiChan <- w2
	}()

	var response models.WeatherApiResponse

	select {
	case w := <-openWeatherChan:
		response.Temp = w.Main.Temp
		response.FeelsLike = w.Main.FeelsLike
		response.Wind = w.Wind.Speed
		response.City = w.Name
		response.Humidity = w.Main.Humidity
		response.Provider = "OpenWeather"
	case w2 := <-weatherApiChan:
		response.Temp = w2.Current.Temp
		response.FeelsLike = w2.Current.FeelsLike
		response.Wind = w2.Current.Wind
		response.City = w2.Location.Name
		response.Humidity = w2.Current.Humidity
		response.Provider = "WeatherApi"
	case <-time.After(10 * time.Second):
		fmt.Println("Timeout - no response received")
		return nil, fmt.Errorf("timeout - no response received")
	}

	weatherData := &models.WeatherData{
		City:         response.City,
		Temp: response.Temp,
		Humidity:     response.Humidity,
		Wind:    response.Wind,
		FeelsLike: response.FeelsLike,
		Provider:     response.Provider,
	}

	err := s.repo.CreateWeather(ctx, weatherData)
	if err != nil {
		fmt.Printf("Warning: failed to save to DB: %v\n", err)
	}

	return &response, nil
}
