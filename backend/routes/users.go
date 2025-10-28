package routes

import (
	"context"
	"net/http"

	"shared-expenses-app/models"

	"github.com/gin-gonic/gin"
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
}
