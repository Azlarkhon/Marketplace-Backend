package controllers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"example.com/m/database"
	"example.com/m/models"
)

func GetProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.DBinstance()

		rows, err := db.Query("SELECT * FROM products")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		defer rows.Close()

		var products []models.Product
		for rows.Next() {
			var product models.Product
			if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Quantity); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}

			products = append(products, product)
		}

		c.JSON(http.StatusOK, products)
	}
}

func GetProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		productName := c.Param("name")

		db := database.DBinstance()

		var product models.Product

		err := db.QueryRow("SELECT * FROM products WHERE name = ?", productName).
			Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Quantity)

		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}

		c.JSON(http.StatusOK, product)
	}
}

func AddProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newProduct models.Product

		db := database.DBinstance()

		if err := c.BindJSON(&newProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		var checkName string

		db.QueryRow("SELECT name FROM products WHERE name = ?", newProduct.Name).
			Scan(&checkName)

		if checkName == newProduct.Name {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The product with equvivalent name already exists"})
			return
		}

		_, err := db.Exec("INSERT INTO products (name, description, price, quantity) VALUES(?, ?, ?, ?)",
			newProduct.Name, newProduct.Description, newProduct.Price, newProduct.Quantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create a product"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product created succesfully"})
	}
}

func UpdateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("product_id")

		var updateProduct models.Product

		if err := c.BindJSON(&updateProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalud input"})
			return
		}

		db := database.DBinstance()

		_, err := db.Exec("UPDATE products SET name = ?, description = ?, price = ?, quantity = ? WHERE product_id = ?",
			updateProduct.Name, updateProduct.Description, updateProduct.Price, updateProduct.Quantity, productID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product updated succesfully"})
	}
}

func DeleteProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("product_id")

		db := database.DBinstance()

		_, err := db.Exec("DELETE FROM products WHERE products.product_id = ?", productID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product deleted sucsesfully"})
	}
}

func FilterProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.DBinstance()

		name := c.Query("name")
		price := c.Query("price")
		quantity := c.Query("quantity")
		sort := c.Query("sort")

		query := "SELECT product_id, name, description, price, quantity FROM products WHERE"
		var args []interface{}
		var conditions []string

		if name != "" {
			conditions = append(conditions, "name LIKE ?")
			args = append(args, "%"+name+"%")
		}
		if price != "" {
			conditions = append(conditions, "price = ?")
			args = append(args, price)
		}
		if quantity != "" {
			conditions = append(conditions, "quantity = ?")
			args = append(args, quantity)
		}

		if len(conditions) == 0 {
			query = "SELECT * FROM products"
		}

		if sort != "" {
			if sort == "asc" {
				query += " ORDER BY name ASC"
			} else if sort == "desc" {
				query += " ORDER BY name DESC"
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sort parameter"})
				return
			}
		}

		query += " " + strings.Join(conditions, " OR ")

		rows, err := db.Query(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query products"})
			return
		}
		defer rows.Close()

		var products []models.Product
		for rows.Next() {
			var product models.Product
			if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Quantity); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan product"})
				return
			}
			products = append(products, product)
		}

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate over products"})
			return
		}

		c.JSON(http.StatusOK, products)
	}
}
