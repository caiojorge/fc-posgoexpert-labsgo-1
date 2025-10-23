package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(weatherHandler *WeatherHandler) *gin.Engine {
	router := gin.Default()

	router.GET("/weather/:zipcode", weatherHandler.GetWeatherByZipCode)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
