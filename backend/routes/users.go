package routes

import (
	"context"
	"net/http"
	"time"

	"shared-expenses-app/models"
	"shared-expenses-app/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserRoutes(router *gin.RouterGroup, pool *pgxpool.Pool) {
	router.GET("list", func(c *gin.Context) {
		rows, err := pool.Query(context.Background(),
			`SELECT user_id, user_name, email, is_guest, password_hash, extract(epoch from created_at)::bigint
			 FROM users ORDER BY created_at DESC`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var users []models.User
		for rows.Next() {
			var u models.User
			err := rows.Scan(&u.UserID, &u.Name, &u.Email, &u.Guest, &u.PasswordHash, &u.CreatedAt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			users = append(users, u)
		}

		c.JSON(http.StatusOK, users)
	})

	// Register a new user
	router.POST("register", func(c *gin.Context) {
		var request struct {
			Name     string `json:"name" binding:"required"`
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password"`
		}

		// Convert request JSON body to struct
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate and convert inputs
		name, err := utils.ValidateName(request.Name)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		email, err := utils.ValidateEmail(request.Email)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		passwordHash, err := utils.HashPassword(request.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// At this point, all inputs are valid

		// Check if user already exists
		var existingID string
		err = pool.QueryRow(
			context.Background(),
			`SELECT user_id FROM users WHERE email = $1`,
			email,
		).Scan(&existingID)

		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}
		if err != pgx.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database query failed"})
			return
		}

		// Insert new user
		var userID string
		err = pool.QueryRow(
			context.Background(),
			`INSERT INTO users (user_name, email, password_hash, created_at)
			 VALUES ($1, $2, $3, $4)
			 RETURNING user_id`,
			name, email, passwordHash, time.Now(),
		).Scan(&userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "user registered successfully",
			"user_id": userID,
		})
	})

	router.POST("login", func(c *gin.Context) {
		var request struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}

		// Convert request JSON body to struct
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		email, err := utils.ValidateEmail(request.Email)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		password, err := utils.HashPassword(request.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var uid, savedPassword string

		err = pool.QueryRow(
			context.Background(),
			`SELECT user_id, password_hash FROM users WHERE email = $1`,
			email,
		).Scan(&uid, &password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}

		if ok := utils.CheckPassword(savedPassword, request.Password); !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}

		// At this point, login is successful

		token, err := utils.GenerateJWT(uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
			return
		}

		// Return JWT token
		c.JSON(http.StatusOK, gin.H{
			"message": "login successful",
			"token":   token,
		})
	})
}
