package route

import (
	"smhome/modules/controller"
	ctrl "smhome/modules/controller"
	mdw "smhome/modules/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoute(r *gin.Engine) {
	r.POST("/api/user/new", ctrl.AddNewUser)
	r.POST("/api/user/login", ctrl.Login)
	r.GET("/api/user/logout", mdw.RequireUser, ctrl.Logout)
	r.GET("/", controller.Public)
}
