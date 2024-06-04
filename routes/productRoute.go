package routes

import (
	controller "example.com/m/controllers"
	"github.com/gin-gonic/gin"
)

func ProductRoutes(incomingRoutes *gin.Engine) {
	// @Summary Get all products
	// @Description Get all products
	// @Tags Products
	// @Success 200 {array} models.Product
	// @Router /products [get]
	incomingRoutes.GET("/products", controller.GetProducts())

	// @Summary Get product by name
	// @Description Get product by name
	// @Tags Products
	// @Param name path string true "Product Name"
	// @Success 200 {object} models.Product
	// @Router /products/{name} [get]
	incomingRoutes.GET("/products/:name", controller.GetProduct())

	// @Summary Filter products
	// @Description Filter products
	// @Tags Products
	// @Success 200 {array} models.Product
	// @Router /productss [get]
	incomingRoutes.GET("/productss", controller.FilterProducts())

	// @Summary Add a new product
	// @Description Add a new product
	// @Tags Products
	// @Param product body models.Product true "Product"
	// @Success 200 {object} models.Product
	// @Router /products [post]
	incomingRoutes.POST("/products", controller.AddProduct())

	// @Summary Update product by ID
	// @Description Update product by ID
	// @Tags Products
	// @Param product_id path string true "Product ID"
	// @Param product body models.Product true "Product"
	// @Success 200 {object} models.Product
	// @Router /products/{product_id} [put]
	incomingRoutes.PUT("/products/:product_id", controller.UpdateProduct())

	// @Summary Delete product by ID
	// @Description Delete product by ID
	// @Tags Products
	// @Param product_id path string true "Product ID"
	// @Success 204 {string} string "No Content"
	// @Router /products/{product_id} [delete]
	incomingRoutes.DELETE("/products/:product_id", controller.DeleteProduct())
}
