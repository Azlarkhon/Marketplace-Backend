package routes

import (
	controller "example.com/m/controllers"
	"github.com/gin-gonic/gin"
)

func CartRoutes(incomingRoutes *gin.Engine) {
	// @Summary Get cart by ID
	// @Description Get cart by ID
	// @Tags Cart
	// @Param cart_id path string true "Cart ID"
	// @Success 200 {object} models.Cart
	// @Router /cart/{cart_id} [get]
	incomingRoutes.GET("/cart/:cart_id", controller.GetCart())

	// @Summary Create a new cart
	// @Description Create a new cart
	// @Tags Cart
	// @Param cart body models.Cart true "Cart"
	// @Success 200 {object} models.Cart
	// @Router /cart [post]
	incomingRoutes.POST("/cart", controller.CreateCart())

	// @Summary Add product to cart
	// @Description Add product to cart
	// @Tags Cart
	// @Param cart_id path string true "Cart ID"
	// @Param product body models.CartProductRequest true "Product"
	// @Success 200 {object} models.CartProduct
	// @Router /cart/{cart_id}/products [post]
	incomingRoutes.POST("/cart/:cart_id/products", controller.AddProductToCart())

	// @Summary Delete product from cart
	// @Description Delete product from cart
	// @Tags Cart
	// @Param cart_product_id path string true "Cart Product ID"
	// @Success 204 {string} string "No Content"
	// @Router /cart_product/{cart_product_id} [delete]
	incomingRoutes.DELETE("/cart_product/:cart_product_id", controller.DeleteProductFromCart())

	// @Summary Update product quantity in cart
	// @Description Update product quantity in cart
	// @Tags Cart
	// @Param cart_id path string true "Cart ID"
	// @Param product body models.CartProductRequest true "Product"
	// @Success 200 {object} models.CartProduct
	// @Router /cart/products [put]
	incomingRoutes.PUT("/cart/products", controller.UpdateQuantityOfProduct())
}
