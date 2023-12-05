package service

import (
	"context"

	pb "weather-service/api/weather/v1"
	"weather-service/internal/biz"
)

type WeatherService struct {
	uc *biz.WeatherUseCase
	pb.UnimplementedWeatherServer
}

func NewWeatherService(uc *biz.WeatherUseCase) *WeatherService {
	return &WeatherService{uc: uc}
}

func (s *WeatherService) GetWeather(ctx context.Context, req *pb.GetWeatherRequest) (*pb.GetWeatherReply, error) {
	weather, err := s.uc.GetWeatherIntorpolate(ctx, float64(req.Lat), float64(req.Lon))
	if err != nil {
		return nil, err
	}
	return &pb.GetWeatherReply{
		Data: &pb.GetWeatherReply_Data{
			Name:        weather.Data.Name,
			Description: weather.Data.Description,
		},
		Temperature: float32(weather.Temperature),
	}, nil
}
