package routes

import (
	controller "example.com/m/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	// @Summary Get all users
	// @Description Get all users
	// @Tags Users
	// @Success 200 {array} models.User
	// @Router /users [get]
	incomingRoutes.GET("/users", controller.GetUsers())

	// @Summary Get user by ID
	// @Description Get user by ID
	// @Tags Users
	// @Param user_id path string true "User ID"
	// @Success 200 {object} models.User
	// @Router /users/{user_id} [get]
	incomingRoutes.GET("/users/:user_id", controller.GetUser())

	// @Summary User signup
	// @Description User signup
	// @Tags Users
	// @Param user body models.User true "User"
	// @Success 200 {object} models.User
	// @Router /users/signup [post]
	incomingRoutes.POST("/users/signup", controller.SignUp())

	// @Summary User login
	// @Description User login
	// @Tags Users
	// @Param user body models.User true "User"
	// @Success 200 {object} models.User
	// @Router /users/login [post]
	incomingRoutes.POST("/users/login", controller.Login())

	// @Summary Update user
	// @Description Update user
	// @Tags Users
	// @Param user_id path string true "User ID"
	// @Param user body models.User true "User"
	// @Success 200 {object} models.User
	// @Router /users/{user_id} [put]
	incomingRoutes.PUT("/users/:user_id", controller.UpdateUser())
}
