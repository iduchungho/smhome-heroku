package controller

import (
	"net/http"
	"os"
	"smhome/modules/model"
	"smhome/modules/service"
	"smhome/modules/utils"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	user, _ := service.NewEntityContext("user")
	_, err := user.FindDocument("username", body.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	hashPass, _ := user.GetElement("password")
	err = utils.ComparePassword(*hashPass, body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Generate a jwt token
	id, _ := user.GetElement("id")
	token := utils.GenerateToken(*id)

	// Sign and get the complete encode token as a string using the secret
	tokenString, errToken := token.SignedString([]byte(os.Getenv("SECRET")))
	if errToken != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	// send it back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
	return
}

func AddNewUser(c *gin.Context) {
	var userMd model.User
	newUser, err := service.NewEntityContext("user")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if newUser.SetElement("type", "user") != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if c.BindJSON(&userMd) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "can't read body request",
		})
		return
	}

	hashPass, err := utils.GenPassword(userMd.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userMd.Password = string(hashPass)

	errIs := newUser.InsertData(userMd)
	if errIs != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errIs.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": userMd,
	})
	return
}

func Logout(c *gin.Context) {
	c.SetCookie("Authorization", "",
		-1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    "your session has been wiped",
	})
}

func Public(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
