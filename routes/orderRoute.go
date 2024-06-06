package routes

import (
	controller "example.com/m/controllers"
	"github.com/gin-gonic/gin"
)

func OrderRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/orders", controller.GetOrders())
	incomingRoutes.GET("/orders/:order_id", controller.GetOrder())
	incomingRoutes.GET("/orders/order_details/:order_id", controller.GetOrderDetail())

	incomingRoutes.POST("/orders", controller.AddOrder())
	incomingRoutes.PUT("/orders/:order_id", controller.UpdateStatus())
	incomingRoutes.DELETE("/orders/:order_id", controller.DeleteOrder())
}
