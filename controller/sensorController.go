package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"smhome/service"
)

func GetTemperature(c *gin.Context) {
	nSensors, err := service.NewEntityContext("sensors")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "can't create sensors models",
		})
		return
	}
	err = nSensors.SetElement("type", "temperature")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	res, errSen := nSensors.GetEntity("")
	if errSen != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errSen.Error(),
		})
		return
	}

	errIs := nSensors.InsertData(res)
	if errIs != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errIs.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": res,
	})
}

func GetHumidity(c *gin.Context) {
	nSensors, err := service.NewEntityContext("sensors")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "can't create sensors models",
		})
		return
	}
	err = nSensors.SetElement("type", "humidity")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	res, errSen := nSensors.GetEntity("")
	if errSen != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errSen.Error(),
		})
		return
	}

	errIs := nSensors.InsertData(res)
	if errIs != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errIs.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": res,
	})
	return
}
