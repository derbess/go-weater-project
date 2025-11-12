package models

type OpenWeather struct {
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
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
		WindKph   float64 `json:"wind_kph"`
	} `json:"current"`
}

type WeatherApiResponse struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feelslike"`
	WindKph   float64 `json:"wind_kph"`
	City      string  `json:"city"`
	Provider  string  `json:"provider"`
}
