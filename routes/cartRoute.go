package routes

import (
	controller "example.com/m/controllers"
	"github.com/gin-gonic/gin"
)

func CartRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.GET("/cart/:cart_id", controller.GetCart())

	incomingRoutes.POST("/cart", controller.CreateCart())
	incomingRoutes.POST("/cart/:cart_id/products", controller.AddProductToCart())

	incomingRoutes.DELETE("/cart_product/:cart_product_id", controller.DeleteProductFromCart())
	incomingRoutes.PUT("/cart/products", controller.UpdateQuantityOfProduct())
}
