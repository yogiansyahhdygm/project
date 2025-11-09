package controllers

import (
	"net/http"

	"toko/models"
	"toko/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetProfil returns pelanggan profile for the logged-in user
func GetProfil(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	claims := c.MustGet("claims").(*utils.Claims)
	userID := claims.UserID

	var pr models.Profil
	if err := db.Where("user_id = ?", userID).First(&pr).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Profil anda tidak ditemukan. Silahkan melakukan Register!"})
		return
	}
	c.JSON(http.StatusOK, pr)
}

// UpdateProfil update pelanggan profile for logged-in user
func UpdateProfil(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	claims := c.MustGet("claims").(*utils.Claims)
	userID := claims.UserID

	var pro models.Profil
	if err := db.Where("user_id = ?", userID).First(&pro).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profil anda tidak ditemukan. Silahkan melakukan Register!"})
		return
	}

	var input struct {
		Nama   string `json:"nama"`
		Email  string `json:"email"`
		Alamat string `json:"alamat"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	updateData := map[string]interface{}{}
	if input.Nama != "" {
		updateData["nama"] = input.Nama
	}
	if input.Email != "" {
		updateData["email"] = input.Email
	}
	if input.Alamat != "" {
		updateData["alamat"] = input.Alamat
	}

	if err := db.Model(&pro).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "data profil berhasil diperbarui",
		"Profil":  pro,
	})

}
