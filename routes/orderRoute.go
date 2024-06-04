package routes

import (
	controller "example.com/m/controllers"
	"github.com/gin-gonic/gin"
)

func OrderRoutes(incomingRoutes *gin.Engine) {
	// @Summary Get all orders
	// @Description Get all orders
	// @Tags Orders
	// @Success 200 {array} models.Order
	// @Router /orders [get]
	incomingRoutes.GET("/orders", controller.GetOrders())

	// @Summary Get order by ID
	// @Description Get order by ID
	// @Tags Orders
	// @Param order_id path string true "Order ID"
	// @Success 200 {object} models.Order
	// @Router /orders/{order_id} [get]
	incomingRoutes.GET("/orders/:order_id", controller.GetOrder())

	// @Summary Get order details by order ID
	// @Description Get order details by order ID
	// @Tags Orders
	// @Param order_id path string true "Order ID"
	// @Success 200 {object} models.Order_details
	// @Router /orders/order_details/{order_id} [get]
	incomingRoutes.GET("/orders/order_details/:order_id", controller.GetOrderDetail())

	// @Summary Add a new order
	// @Description Add a new order
	// @Tags Orders
	// @Param order body models.Order true "Order"
	// @Success 200 {object} models.Order
	// @Router /orders [post]
	incomingRoutes.POST("/orders", controller.AddOrder())

	// @Summary Update order status by ID
	// @Description Update order status by ID
	// @Tags Orders
	// @Param order_id path string true "Order ID"
	// @Param status body models.OrderStatus true "Order Status"
	// @Success 200 {object} models.Order
	// @Router /orders/{order_id} [put]
	incomingRoutes.PUT("/orders/:order_id", controller.UpdateStatus())

	// @Summary Delete order by ID
	// @Description Delete order by ID
	// @Tags Orders
	// @Param order_id path string true "Order ID"
	// @Success 204 {string} string "No Content"
	// @Router /orders/{order_id} [delete]
	incomingRoutes.DELETE("/orders/:order_id", controller.DeleteOrder())
}
