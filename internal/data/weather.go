package data

import (
	"context"
	"time"
	"weather-service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type Weather struct {
	Id          uint64 `gorm:"primaryKey"`
	Temperature float64
	WeahterID   uint
	Time        time.Time
}

func (m Weather) modelToResponse() *biz.Weather {
	return &biz.Weather{
		Id:          m.Id,
		Temperature: m.Temperature,
		WeahterID:   m.WeahterID,
		Time:        m.Time,
	}
}

type weatherRepo struct {
	data   *Data
	logger *log.Helper
}

func NewWeatherRepo(data *Data, logger log.Logger) biz.WeatherRepo {
	return &weatherRepo{data: data, logger: log.NewHelper(logger)}
}

// Create implements biz.WeatherRepo.
func (r *weatherRepo) Create(ctx context.Context, weathers []*biz.Weather) error {
	weathersDB := make([]Weather, 0)
	for _, weather := range weathers {
		weatherDB := Weather{}
		weatherDB.Temperature = weather.Temperature
		weatherDB.WeahterID = weather.WeahterID
		weatherDB.Time = time.Now()
		weathersDB = append(weathersDB, weatherDB)
	}
	if err := r.data.db.Create(&weathersDB).Error; err != nil {
		return err
	}
	return nil
}

// GetWeather implements biz.WeatherRepo.
func (r *weatherRepo) GetWeather(context.Context) ([]*biz.Weather, error) {
	var weathersDB []Weather
	result := r.data.db.Model(&Weather{}).Order("time DESC").Limit(40).Find(&weathersDB)
	if result.Error != nil {
		return nil, result.Error
	}

	weathers := make([]*biz.Weather, 0)
	for _, weather := range weathersDB {
		weathers = append(weathers, weather.modelToResponse())
	}
	return weathers, nil
}
