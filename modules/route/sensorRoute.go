package route

import (
	"github.com/gin-gonic/gin"
	"smhome/modules/controller"
)

func SenSorRoute(r *gin.Engine) {
	r.GET("/api/sensor/temperature", controller.GetTemperature)
	r.GET("/api/sensor/humidity", controller.GetHumidity)
}
