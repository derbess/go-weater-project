package models

import (
	"time"
)
type OpenWeather struct {
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		Humidity  float64 `json:"humidity"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Name string `json:"name"`
}

type WeatherApi struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
	Current struct {
		Temp      float64 `json:"temp_c"`
		FeelsLike float64 `json:"feelslike_c"`
		Wind      float64 `json:"wind_kph"`
		Humidity  float64 `json:"humidity"`
	} `json:"current"`
}

type WeatherApiResponse struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feelslike"`
	Wind      float64 `json:"wind"`
	City      string  `json:"city"`
	Humidity  float64 `json:"humidity"`
	Provider  string  `json:"provider"`
}

type WeatherData struct {
	ID        int64   `json:"id"`
	City      string  `json:"city"`
	Temp      float64 `json:"temp"`
	Humidity  float64 `json:"humidity"`
	Wind      float64 `json:"wind"`  
	FeelsLike float64 `json:"feelslike"`
	Provider  string  `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
