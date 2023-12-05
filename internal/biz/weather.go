package biz

import (
	"context"
	"time"
	"weather-service/internal/utils"

	"github.com/go-kratos/kratos/v2/log"
)

type Weather struct {
	Id          uint64
	Temperature float64
	WeahterID   uint
	Time        time.Time
}

type WeatherRepo interface {
	Create(context.Context, []*Weather) error
	GetWeather(context.Context) ([]*Weather, error)
}

type WeatherUseCase struct {
	repo   WeatherRepo
	logger *log.Helper
}

func NewWeatherUseCase(repo WeatherRepo, logger log.Logger) *WeatherUseCase {
	return &WeatherUseCase{repo: repo, logger: log.NewHelper(logger)}
}

func (uc *WeatherUseCase) Create(ctx context.Context, weather []*Weather) error {
	return uc.repo.Create(ctx, weather)
}

func (uc *WeatherUseCase) GetWeatherIntorpolate(ctx context.Context, lat, lon float64) (*WeaterResponse, error) {
	weathers, err := uc.repo.GetWeather(ctx)
	if err != nil {
		return nil, err
	}
	var coords []utils.Coordinate
	for i := 0; i < len(utils.Latitudes); i++ {
		coords = append(coords, utils.Coordinate{Latitude: utils.Latitudes[i], Longitude: utils.Longitudes[i]})
	}

	tempGrid := make([][]float64, len(utils.Latitudes))
	for i := range tempGrid {
		tempGrid[i] = make([]float64, len(utils.Longitudes))
	}
	for i, lat := range utils.Latitudes {
		for j, lon := range utils.Longitudes {
			index := utils.FindIndexForCoord(lat, lon, coords)
			if index != -1 {
				tempGrid[i][j] = weathers[index].Temperature
			}
		}
	}

	interpolatedTemp, index := utils.InterpolateTemperature(lat, lon, utils.Latitudes, utils.Longitudes, tempGrid)
	return &WeaterResponse{
		Data:        indexMap[uint(index)],
		Temperature: interpolatedTemp,
	}, nil
}

type WeaterResponse struct {
	Data        WeatherData `json:"weather"`
	Temperature float64     `json:"temp"`
}

type WeatherData struct {
	Name        string
	Description string
}

var indexMap = map[uint]WeatherData{
	200: {
		Name:        "Гроза",
		Description: "Гроза с небольшим дождем",
	},
	201: WeatherData{
		Name:        "Гроза",
		Description: "Гроза с дождем",
	},
	202: WeatherData{
		Name:        "Гроза",
		Description: "Гроза с сильным дождем",
	},
	210: WeatherData{
		Name:        "Гроза",
		Description: "Легкая гроза",
	},
	211: WeatherData{
		Name:        "Гроза",
		Description: "Гроза",
	},
	212: WeatherData{
		Name:        "Гроза",
		Description: "Сильная гроза",
	},
	221: WeatherData{
		Name:        "Гроза",
		Description: "Неровная гроза",
	},
	230: WeatherData{
		Name:        "Гроза",
		Description: "Гроза с небольшим моросью",
	},
	231: WeatherData{
		Name:        "Гроза",
		Description: "Гроза с моросью",
	},
	232: WeatherData{
		Name:        "Гроза",
		Description: "Гроза с сильным моросящим дождем",
	},
	300: WeatherData{
		Name:        "Моросить",
		Description: "Слабый дождь",
	},
	301: WeatherData{
		Name:        "Моросить",
		Description: "Моросить",
	},
	302: WeatherData{
		Name:        "Моросить",
		Description: "Сильный дождь",
	},
	310: WeatherData{
		Name:        "Моросить",
		Description: "Небольшой дождь",
	},
	311: WeatherData{
		Name:        "Моросить",
		Description: "Моросящий дождь",
	},
	312: WeatherData{
		Name:        "Моросить",
		Description: "Сильный моросящий дождь",
	},
	313: WeatherData{
		Name:        "Моросить",
		Description: "Дождь с моросью",
	},
	314: WeatherData{
		Name:        "Моросить",
		Description: "Сильный ливень и изморость",
	},
	321: WeatherData{
		Name:        "Моросить",
		Description: "Дождь из душа",
	},
	500: WeatherData{
		Name:        "Дождь",
		Description: "Легкий дождь",
	},
	501: WeatherData{
		Name:        "Дождь",
		Description: "Умеренный дождь",
	},
	502: WeatherData{
		Name:        "Дождь",
		Description: "Сильный дождь",
	},
	503: WeatherData{
		Name:        "Дождь",
		Description: "Очень сильный дождь",
	},
	504: WeatherData{
		Name:        "Дождь",
		Description: "Сильный дождь",
	},
	511: WeatherData{
		Name:        "Дождь",
		Description: "Ледяной дождь",
	},
	520: WeatherData{
		Name:        "Дождь",
		Description: "Легкий дождь",
	},
	521: WeatherData{
		Name:        "Дождь",
		Description: "Ливень",
	},
	522: WeatherData{
		Name:        "Дождь",
		Description: "Сильный ливень",
	},
	531: WeatherData{
		Name:        "Дождь",
		Description: "Неровный дождь",
	},
	600: WeatherData{
		Name:        "Снег",
		Description: "Легкий снег",
	},
	601: WeatherData{
		Name:        "Снег",
		Description: "Снег",
	},
	602: WeatherData{
		Name:        "Снег",
		Description: "Сильный снегопад",
	},
	611: WeatherData{
		Name:        "Снег",
		Description: "Мокрый снег",
	},
	612: WeatherData{
		Name:        "Снег",
		Description: "Легкий дожль с мокрым снегом",
	},
	613: WeatherData{
		Name:        "Снег",
		Description: "Мокрый дождь",
	},
	615: WeatherData{
		Name:        "Снег",
		Description: "Небольшой дождь и снег",
	},
	616: WeatherData{
		Name:        "Снег",
		Description: "Дождь и снег",
	},
	620: WeatherData{
		Name:        "Снег",
		Description: "Легкий ливневый снег",
	},
	621: WeatherData{
		Name:        "Снег",
		Description: "Ливневый снег",
	},
	622: WeatherData{
		Name:        "Снег",
		Description: "Сильный снегопад",
	},
	701: WeatherData{
		Name:        "Туман",
		Description: "Туман",
	},
	711: WeatherData{
		Name:        "Смог",
		Description: "Смог",
	},
	721: WeatherData{
		Name:        "Туман",
		Description: "Туман",
	},
	731: WeatherData{
		Name:        "Пыль",
		Description: "Пыль",
	},
	741: WeatherData{
		Name:        "Туман",
		Description: "Туман",
	},
	751: WeatherData{
		Name:        "Песчанная буря",
		Description: "Песчанная буря",
	},
	761: WeatherData{
		Name:        "Пыль",
		Description: "Пыль",
	},
	762: WeatherData{
		Name:        "Пепель",
		Description: "Пепель",
	},
	771: WeatherData{
		Name:        "Шквалы",
		Description: "Шквалы",
	},
	781: WeatherData{
		Name:        "Торнадо",
		Description: "Торнадо",
	},
	800: WeatherData{
		Name:        "Чисто",
		Description: "Все хорошо",
	},
	801: WeatherData{
		Name:        "Облачно",
		Description: "Мало облаков: 11-25%",
	},
	802: WeatherData{
		Name:        "Облачно",
		Description: "Рассеянные облака: 25-50%",
	},
	803: WeatherData{
		Name:        "Облачно",
		Description: "Разорванные облака: 51-84%",
	},
	804: WeatherData{
		Name:        "Облачно",
		Description: "Пасмурная облачность: 85-100%",
	},
}
