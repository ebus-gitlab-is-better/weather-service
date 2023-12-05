package server

import (
	"context"
	"fmt"
	"time"
	"weather-service/internal/biz"
	"weather-service/internal/conf"
	"weather-service/internal/utils"

	"github.com/go-co-op/gocron"
)

type Parser struct {
}

func NewParser(c *conf.Data, uc *biz.WeatherUseCase) *Parser {
	loc, _ := time.LoadLocation("Europe/Moscow")
	cron := gocron.NewScheduler(loc)
	cron.Every("1m").Do(func() {
		var coords []utils.Coordinate
		for i := 0; i < len(utils.Latitudes); i++ {
			coords = append(coords, utils.Coordinate{Latitude: utils.Latitudes[i], Longitude: utils.Longitudes[i]})
		}
		temperatures, arrayIndex, err := utils.FetchTemperaturesForCoordinates(coords, c.Api)
		if err != nil {
			fmt.Println("Error fetching temperatures:", err)
		}
		weathers := make([]*biz.Weather, 0)
		for i := 0; i < len(temperatures); i++ {
			weathers = append(weathers, &biz.Weather{
				Temperature: temperatures[i],
				WeahterID:   arrayIndex[i],
			})
		}
		uc.Create(context.TODO(), weathers)
	})
	cron.StartAsync()
	return &Parser{}
}

func (s *Parser) Start(ctx context.Context) error {
	return nil
}

func (s *Parser) Stop(ctx context.Context) error {
	return nil
}
