package services

import (
	"context"
	"fmt"
	"go-first-project/client"
	"go-first-project/models"
	"go-first-project/repository"
	"time"
)

type WeatcherService struct {
	repo   *repository.WeatherRepository
	client *client.Client
}

func NewWeatcherService(repo *repository.WeatherRepository, client *client.Client) *WeatcherService {
	return &WeatcherService{
		repo:   repo,
		client: client,
	}
}

func (s *WeatcherService) GetWeather(ctx context.Context, city string) (*models.WeatherApiResponse, error) {
	start := time.Now()

	apiCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	openWeatherChan := make(chan *models.OpenWeather)
	weatherApiChan := make(chan *models.WeatherApi)

	go func() {
		w, err := s.client.FetchOpenWeather(apiCtx, city)
		if err != nil {
			fmt.Printf("Error fetching OpenWeather: %v\n", err)
			return
		}
		fmt.Printf("OpenWeather time: %v\n", time.Since(start))
		openWeatherChan <- w
	}()

	go func() {
		w2, err := s.client.FetchWeatherApi(apiCtx, city)
		if err != nil {
			fmt.Printf("Error fetching WeatherApi: %v\n", err)
			return
		}
		fmt.Printf("WeatherApi time: %v\n", time.Since(start))
		weatherApiChan <- w2
	}()

	var response models.WeatherApiResponse

	select {
	case w := <-openWeatherChan:
		cancel()
		response.Temp = w.Main.Temp
		response.FeelsLike = w.Main.FeelsLike
		response.Wind = w.Wind.Speed
		response.City = w.Name
		response.Humidity = w.Main.Humidity
		response.Provider = "OpenWeather"
	case w2 := <-weatherApiChan:
		cancel()
		response.Temp = w2.Current.Temp
		response.FeelsLike = w2.Current.FeelsLike
		response.Wind = w2.Current.Wind
		response.City = w2.Location.Name
		response.Humidity = w2.Current.Humidity
		response.Provider = "WeatherApi"
	case <-apiCtx.Done():
		fmt.Printf("Timeout - no response received")
		return nil, fmt.Errorf("timeout - no response received")
	}

	weatherData := &models.WeatherData{
		City:      response.City,
		Temp:      response.Temp,
		Humidity:  response.Humidity,
		Wind:      response.Wind,
		FeelsLike: response.FeelsLike,
		Provider:  response.Provider,
	}

	err := s.repo.CreateWeather(ctx, weatherData)
	if err != nil {
		fmt.Printf("Warning: failed to save to DB: %v\n", err)
	}

	return &response, nil
}
