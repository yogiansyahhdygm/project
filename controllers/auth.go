package controllers

import (
	"net/http"

	"toko/models"
	"toko/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func Login(c *gin.Context) {
	var input AuthRequest
	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau Password Salah!"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "password salah"})
		return
	}

	token, err := utils.GenerateToken(int(user.ID), user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"role":  user.Role,
	})
}

func Register(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var req AuthRequest

	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username & password required"})
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengenkripsi password"})
		return
	}

	rl := req.Role
	if rl == "" {
		rl = "pelanggan"
	}

	user := models.User{
		Username: req.Username,
		Password: string(hashed),
		Role:     rl,
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// kita buat profil
	profil := models.Profil{
		UserID: &user.ID,
	}

	if err := db.Create(&profil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat profil"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "registrasi berhasil",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
		"profi": gin.H{
			"id":      profil.ID,
			"user_id": profil.UserID,
		},
	})

}
