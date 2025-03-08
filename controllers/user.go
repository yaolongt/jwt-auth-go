package controllers

import (
	"encoding/json"
	"go-jwt/configs"
	"go-jwt/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserBodyRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserController struct {
	DB *gorm.DB
}

func NewUserController() *UserController {
	db := configs.GetDBInstance().DB
	return &UserController{DB: db}
}

func (t *UserController) SignUp(c *gin.Context) {
	var body UserBodyRequest

	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body; Invalid body request",
		})
		return
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	user := models.User{Email: body.Email, Password: string(hash)}
	result := t.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func (t *UserController) Login(c *gin.Context) {
	var body UserBodyRequest

	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body; Invalid body request",
		})
		return
	}

	var user models.User
	t.DB.First(&user, "email=?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	ttl := time.Hour * 3
	now := time.Now().UTC()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub":   user.ID,
		"exp":   now.Add(ttl).Unix(), // Expiration time
		"iat":   now.Unix(),          // Time issued
		"nbf":   now.Unix(),          // Time before which is invalid
		"email": user.Email,
	})

	key, err := jwt.ParseRSAPrivateKeyFromPEM(configs.PRIV)
	tokenString, err := token.SignedString(key)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		log.Fatal(err)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*3, "", "", true, true)
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (t *UserController) Validate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged in and valid",
	})
}
