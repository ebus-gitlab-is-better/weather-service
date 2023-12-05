package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Coordinate struct {
	Latitude  float64
	Longitude float64
}

type WeatherData struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weater []struct {
		Id uint `json:"id"`
	} `json:"weather"`
}

func InterpolateTemperature(lat, lon float64, latGrid, lonGrid []float64, tempGrid [][]float64) (float64, int) {
	latIndex := findNearestIndex(lat, latGrid)
	lonIndex := findNearestIndex(lon, lonGrid)

	x1, x2 := latGrid[latIndex], latGrid[latIndex+1]
	y1, y2 := lonGrid[lonIndex], lonGrid[lonIndex+1]

	Q11, Q21 := tempGrid[latIndex][lonIndex], tempGrid[latIndex+1][lonIndex]
	Q12, Q22 := tempGrid[latIndex][lonIndex+1], tempGrid[latIndex+1][lonIndex+1]

	interpolatedTemp := Q11*(x2-lat)*(y2-lon)/((x2-x1)*(y2-y1)) +
		Q21*(lat-x1)*(y2-lon)/((x2-x1)*(y2-y1)) +
		Q12*(x2-lat)*(lon-y1)/((x2-x1)*(y2-y1)) +
		Q22*(lat-x1)*(lon-y1)/((x2-x1)*(y2-y1))

	return interpolatedTemp, latIndex
}

func findNearestIndex(value float64, grid []float64) int {
	for i := 0; i < len(grid)-1; i++ {
		if value >= grid[i] && value <= grid[i+1] {
			return i
		}
	}
	return -1
}

func fetchTemperature(coord Coordinate, apiKey string) (float64, uint, error) {
	// Construct the API request URL
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric", coord.Latitude, coord.Longitude, apiKey)

	// Create a new HTTP client with a timeout
	client := &http.Client{}

	// Make the HTTP request
	resp, err := client.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}

	// Unmarshal the JSON data
	var data WeatherData
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, 0, err
	}

	// Return the temperature
	return data.Main.Temp, data.Weater[0].Id, nil
}

func FetchTemperaturesForCoordinates(coords []Coordinate, apiKey string) ([]float64, []uint, error) {
	var temperatures []float64
	var arrayIndex []uint
	for _, coord := range coords {
		temp, index, err := fetchTemperature(coord, apiKey)
		if err != nil {
			return nil, nil, err
		}
		arrayIndex = append(arrayIndex, index)
		temperatures = append(temperatures, temp)
	}
	return temperatures, arrayIndex, nil
}

func FindIndexForCoord(lat, lon float64, coords []Coordinate) int {
	for i, coord := range coords {
		if coord.Latitude == lat && coord.Longitude == lon {
			return i
		}
	}
	return -1
}
