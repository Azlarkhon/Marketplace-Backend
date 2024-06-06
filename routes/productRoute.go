package routes

import (
	controller "example.com/m/controllers"
	"github.com/gin-gonic/gin"
)

func ProductRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/products", controller.GetProducts())
	incomingRoutes.GET("/products/:name", controller.GetProduct())
	incomingRoutes.GET("/productss", controller.FilterProducts())

	incomingRoutes.POST("/products", controller.AddProduct())

	incomingRoutes.PUT("/products/:product_id", controller.UpdateProduct())

	incomingRoutes.DELETE("/products/:product_id", controller.DeleteProduct())
}
