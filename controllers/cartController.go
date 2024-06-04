package controllers

import (
	"database/sql"
	"net/http"

	"example.com/m/database"
	"example.com/m/models"
	"github.com/gin-gonic/gin"
)

func CreateCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cartRequest models.CartRequest

		if err := c.BindJSON(&cartRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		db := database.DBinstance()

		result, err := db.Exec("INSERT INTO cart (user_id) VALUES (?)", cartRequest.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
			return
		}

		cartID, err := result.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart ID"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Cart created successfully", "cart_id": cartID})
	}
}

func GetCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		cartID := c.Param("cart_id")
		db := database.DBinstance()

		var cart models.Cart
		err := db.QueryRow("SELECT cart_id, user_id FROM cart WHERE cart_id = ?", cartID).Scan(&cart.CartID, &cart.UserID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rows, err := db.Query("SELECT cart_product_id, cart_id, product_id, quantity FROM cart_products WHERE cart_id = ?", cartID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var cartProducts []models.CartProduct
		for rows.Next() {
			var cartProduct models.CartProduct
			if err := rows.Scan(&cartProduct.CartProductID, &cartProduct.CartID, &cartProduct.ProductID, &cartProduct.Quantity); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			cartProducts = append(cartProducts, cartProduct)
		}

		c.JSON(http.StatusOK, gin.H{"cart": cart, "products": cartProducts})
	}
}

func AddProductToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		cartID := c.Param("cart_id")
		var cartProducts []models.CartProductRequest

		if err := c.BindJSON(&cartProducts); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		db := database.DBinstance()
		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
			return
		}

		defer func() {
			if err != nil {
				tx.Rollback()
				return
			}
		}()

		for _, cartProduct := range cartProducts {
			_, err = tx.Exec("INSERT INTO cart_products (cart_id, product_id, quantity) VALUES (?, ?, ?)",
				cartID, cartProduct.ProductID, cartProduct.Quantity)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product to cart"})
				return
			}
		}

		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Products added to cart successfully"})
	}
}

func DeleteProductFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		cartProductID := c.Param("cart_product_id")
		db := database.DBinstance()

		_, err := db.Exec("DELETE FROM cart_products WHERE cart_product_id = ?", cartProductID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cart product"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Cart product deleted successfully"})
	}
}

func UpdateQuantityOfProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			CartProductID int  `json:"cart_product_id"`
			CartID        int  `json:"cart_id"`
			Increment     bool `json:"increment"` // true to increase, false to decrease
		}

		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		db := database.DBinstance()

		var query string
		if input.Increment {
			query = "UPDATE cart_products SET quantity = quantity + 1 WHERE cart_product_id = ? AND cart_id = ?"
		} else {
			query = "UPDATE cart_products SET quantity = quantity - 1 WHERE cart_product_id = ? AND cart_id = ?"
		}

		result, err := db.Exec(query, input.CartProductID, input.CartID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart product"})
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get affected rows"})
			return
		}

		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart product not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Cart product updated successfully"})
	}
}
