package utils

import (
	"fmt"
	"strconv"
	"strings"
)

var Latitudes = []float64{55.3,
	55.5,
	55.7,
	55.9,
	56.1,
	55.3,
	55.5,
	55.7,
	55.9,
	56.1,
	55.3,
	55.5,
	55.7,
	55.9,
	56.1,
	55.3,
	55.5,
	55.7,
	55.9,
	56.1,
	55.3,
	55.5,
	55.7,
	55.9,
	56.1,
	55.3,
	55.5,
	55.7,
	55.9,
	56.1,
	55.3,
	55.5,
	55.7,
	55.9,
	56.1,
	55.3,
	55.5,
	55.7,
	55.9,
	56.1}
var Longitudes = []float64{36.9,
	36.9,
	36.9,
	36.9,
	36.9,
	37.1,
	37.1,
	37.1,
	37.1,
	37.1,
	37.3,
	37.3,
	37.3,
	37.3,
	37.3,
	37.5,
	37.5,
	37.5,
	37.5,
	37.5,
	37.7,
	37.7,
	37.7,
	37.7,
	37.7,
	37.9,
	37.9,
	37.9,
	37.9,
	37.9,
	38.1,
	38.1,
	38.1,
	38.1,
	38.1,
	38.3,
	38.3,
	38.3,
	38.3,
	38.3}

func ParseCoordinates(input string) ([][2]float64, error) {
	var points [][2]float64
	coords := strings.Split(input, ";")
	for _, coord := range coords {
		parts := strings.Split(coord, ",")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid coordinate format")
		}
		lon, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return nil, err
		}
		lat, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return nil, err
		}
		points = append(points, [2]float64{lon, lat})
	}
	return points, nil
}
