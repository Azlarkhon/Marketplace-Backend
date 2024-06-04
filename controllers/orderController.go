package controllers

import (
	"database/sql"
	"net/http"
	"time"

	"example.com/m/database"
	"example.com/m/models"
	"github.com/gin-gonic/gin"
)

func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.DBinstance()

		var createdAt string

		rows, err := db.Query("SELECT order_id, user_id, total_amount, status, created_at FROM orders")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var orders []models.OrderResponse
		for rows.Next() {
			var order models.OrderResponse
			if err := rows.Scan(&order.Order_id, &order.User_id, &order.Total_amount, &order.Status, &createdAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			createdAtTime, err := time.Parse("2006-01-02 15:04:05", createdAt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			order.Created_at = createdAtTime

			orders = append(orders, order)
		}

		c.JSON(http.StatusOK, orders)
	}
}

func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param("order_id")

		db := database.DBinstance()

		var createdAt string

		var order models.OrderResponse

		err := db.QueryRow("SELECT order_id, user_id, total_amount, status, created_at FROM orders WHERE order_id = ?", orderID).
			Scan(&order.Order_id, &order.User_id, &order.Total_amount, &order.Status, &createdAt)

		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
				return
			}
		}

		createdAtTime, err := time.Parse("2006-01-02 15:04:05", createdAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		order.Created_at = createdAtTime

		c.JSON(http.StatusOK, order)
	}
}

func GetOrderDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.DBinstance()

		orderID := c.Param("order_id")

		rows, err := db.Query("SELECT * FROM order_details WHERE order_id = ?", orderID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var orderDetails []models.Order_details
		for rows.Next() {
			var orderDetail models.Order_details
			err := rows.Scan(&orderDetail.OrderDetailID, &orderDetail.Order_id, &orderDetail.Product_id, &orderDetail.Quantity, &orderDetail.Price)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			orderDetails = append(orderDetails, orderDetail)
		}

		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(orderDetails) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		c.JSON(http.StatusOK, orderDetails)
	}
}

func AddOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newOrder models.Order
		newOrder.Status = "active"
		newOrder.Total_amount = 0

		if err := c.BindJSON(&newOrder); err != nil {
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

		result, err := tx.Exec("INSERT INTO orders (user_id, total_amount, status) VALUES (?, ?, ?)",
			newOrder.User_id, newOrder.Total_amount, newOrder.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create an order"})
			return
		}

		orderID, err := result.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order ID"})
			return
		}

		for _, detail := range newOrder.OrderDetails {
			var availableQuantity int
			err := tx.QueryRow("SELECT quantity FROM products WHERE product_id = ?", detail.Product_id).Scan(&availableQuantity)
			if err != nil {
				if err == sql.ErrNoRows {
					c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			if availableQuantity < detail.Quantity {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient product quantity of"})
				return
			}

			price := GetProductPrice(c, float64(detail.Product_id), detail.Quantity)
			if price == -1 {
				return
			}

			_, err = tx.Exec("INSERT INTO order_details (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)",
				orderID, detail.Product_id, detail.Quantity, price)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order details"})
				return
			}

			_, err = tx.Exec("UPDATE products SET quantity = quantity - ? WHERE product_id = ?",
				detail.Quantity, detail.Product_id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product quantity"})
				return
			}

			newOrder.Total_amount += price
		}

		_, err = tx.Exec("UPDATE orders SET total_amount = ? WHERE order_id = ? AND user_id = ?",
			newOrder.Total_amount, orderID, newOrder.User_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order total amount"})
			return
		}

		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order created successfully"})
	}
}

func UpdateStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param("order_id")

		var status string

		db := database.DBinstance()

		db.QueryRow("SELECT status FROM orders WHERE order_id = ?", orderID).Scan(&status)

		if status == "active" {
			status = "inactive"
			_, err := db.Exec("UPDATE orders SET status = ? WHERE order_id = ?", status, orderID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}
		} else {
			status = "active"
			_, err := db.Exec("UPDATE orders SET status = ? WHERE order_id = ?", status, orderID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order updated succesfully"})
	}
}

func DeleteOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param("order_id")

		db := database.DBinstance()

		_, err := db.Exec("DELETE FROM orders WHERE orders.order_id = ?", orderID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order deleted sucsesfully"})
	}
}

func GetProductPrice(c *gin.Context, id float64, quantity int) float64 {
	db := database.DBinstance()

	var price float64

	err := db.QueryRow("SELECT price FROM products WHERE product_id = ?", id).Scan(&price)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return -1
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return -1
	}
	return price * float64(quantity)
}

func GetTotalAmount(c *gin.Context, orderID int, userID int) float64 {
	db := database.DBinstance()

	var totalAmount float64

	rows, err := db.Query("SELECT product_id, quantity FROM order_details WHERE order_id = ?", orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return -1
	}
	defer rows.Close()

	for rows.Next() {
		var productID int
		var quantity int

		err := rows.Scan(&productID, &quantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return -1
		}

		price := GetProductPrice(c, float64(productID), quantity)
		if price == -1 {
			return -1
		}

		totalAmount += price
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return -1
	}

	return totalAmount
}
