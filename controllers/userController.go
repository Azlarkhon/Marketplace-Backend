package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"example.com/m/database"
	"example.com/m/models"
)

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.DBinstance()

		var createdAtStr string

		rows, err := db.Query("SELECT * FROM users")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var users []models.User
		for rows.Next() {
			var user models.User
			if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Address, &user.Password, &user.Age, &createdAtStr); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			user.Created_at = createdAt

			users = append(users, user)
		}

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		userID := c.Param("user_id")

		db := database.DBinstance()

		var user models.User
		var createdAtStr string

		err := db.QueryRow("SELECT * FROM users WHERE user_id = ?", userID).
			Scan(&user.ID, &user.Username, &user.Email, &user.Address, &user.Password, &user.Age, &createdAtStr)

		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		user.Created_at = createdAt

		c.JSON(http.StatusOK, user)
	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newUser models.User

		if err := c.BindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		db := database.DBinstance()
		var checkEmail string

		db.QueryRow("SELECT email FROM users WHERE email = ?", newUser.Email).
			Scan(&checkEmail)

		if checkEmail == newUser.Email {
			log.Printf("Email is already added")
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		newUser.Password = string(hashedPassword)

		_, err = db.Exec("INSERT INTO users (username, email, address, password, age) VALUES (?, ?, ?, ?, ?)",
			newUser.Username, newUser.Email, newUser.Address, newUser.Password, newUser.Age)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.BindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		db := database.DBinstance()

		var user models.User
		var storedPassword string
		var createdAtStr string

		err := db.QueryRow("SELECT * FROM users WHERE email = ?", loginRequest.Email).
			Scan(&user.ID, &user.Username, &user.Email, &user.Address, &storedPassword, &user.Age, &createdAtStr)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query user"})
			return
		}

		createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		user.Created_at = createdAt

		if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(loginRequest.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
	}
}

func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		userID := c.Param("user_id")

		var updateUser models.User

		if err := c.BindJSON(&updateUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		db := database.DBinstance()

		var storedPassword string

		err := db.QueryRow("SELECT password FROM users WHERE user_id = ?", userID).Scan(&storedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			log.Printf("Error querying user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query user"})
			return
		}

		if updateUser.Password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateUser.Password), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
				return
			}
			updateUser.Password = string(hashedPassword)
		} else {
			updateUser.Password = storedPassword
		}

		_, err = db.Exec("UPDATE users SET username = ?, email = ?, address = ?, password = ?, age = ? WHERE user_id = ?",
			updateUser.Username, updateUser.Email, updateUser.Address, updateUser.Password, updateUser.Age, userID)
		if err != nil {
			log.Printf("Error updating user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
	}
}
