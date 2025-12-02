package models

type WeatherQueryParams struct {
	City string `validate:"required,min=2,max=100"`
}
