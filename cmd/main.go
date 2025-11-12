package main

import (
	"go-first-project/services"
	"go-first-project/handlers"
	"net/http"
)

func main() {

	wServices := services.WeatcherService{}
 	wHandlers := handlers.WeatherHandler{WeatcherService: wServices}


	http.HandleFunc("/weather", wHandlers.GetWeatherHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}

	
}
