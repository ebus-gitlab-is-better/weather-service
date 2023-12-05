package route

import (
	"context"
	"net/http"
	"weather-service/internal/biz"
	"weather-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type WeatherRoute struct {
	uc *biz.WeatherUseCase
}

func NewWeatherRoute(uc *biz.WeatherUseCase) *WeatherRoute {
	return &WeatherRoute{uc: uc}
}

func (r *WeatherRoute) Register(router *gin.RouterGroup) {
	router.GET("/:coordinates", r.getWeather)
}

// @Summary	Get Weather
// @Accept		json
// @Produce	json
// @Tags		map
// @Param			coordinate	path	string	true	"[{longitude},{latitude}]"
//
//	@Success	200	{object}	biz.WeaterResponse
//
// @Failure	401
// @Failure	403
// @Failure	500
// @Failure	400
// @Failure	404
// @Router		/weather/{coordinates}/ [get]
func (r *WeatherRoute) getWeather(c *gin.Context) {
	coordinates := c.Param("coordinates")
	points, err := utils.ParseCoordinates(coordinates)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	weather, err := r.uc.GetWeatherIntorpolate(context.TODO(), points[0][1], points[0][0])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, weather)
}
