package main

import (
	"fmt"
	"github.com/shawntoffel/darksky"
)

// print current temp in new york
func main() {

	client := darksky.New("Api key")

	request := darksky.ForecastRequest{}
	request.Latitude = 40.7128
	request.Longitude = -74.0059

	response, err := client.Forecast(request)

	if err != nil {
		fmt.Println(err.Error())

		return
	}

	fmt.Println(response.Currently.Temperature)
}
