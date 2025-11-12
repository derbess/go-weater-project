package services

import (
	"context"
	"fmt"
	"go-first-project/internal"
	"go-first-project/models"
	"net/http"
	"time"
)

type WeatcherService struct {
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
		response.WindKph = w.Wind.Speed
		response.City = w.Name
		response.Provider = "OpenWeather"
	case w2 := <-weatherApiChan:
		response.Temp = w2.Current.Temp
		response.FeelsLike = w2.Current.FeelsLike
		response.WindKph = w2.Current.WindKph
		response.City = w2.Location.Name
		response.Provider = "WeatherApi"
	case <-time.After(10 * time.Second):
		fmt.Println("Timeout - no response received")
		return nil, fmt.Errorf("timeout - no response received")
	}

	return &response, nil
}
