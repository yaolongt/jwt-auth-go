package middleware

import (
	"fmt"
	"go-jwt/configs"
	"go-jwt/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type Middleware struct {
	DB *gorm.DB
}

func NewMiddleware() *Middleware {
	db := configs.GetDBInstance().DB
	return &Middleware{DB: db}
}

func (t *Middleware) Auth(c *gin.Context) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(configs.PUB)

	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Unable to authorize user",
		})
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return key, nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Token invalid",
		})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token invalid",
			})
			return
		}

		var user models.User
		t.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unable to authorize user",
			})
			return
		}

		c.Set("user", user)

		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
