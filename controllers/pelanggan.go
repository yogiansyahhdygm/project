package controllers

import (
	"net/http"

	"toko/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GET all /pelanggan
func GetAllPelanggan(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var pelanggan []models.Profil
	if err := db.Preload("User").Find(&pelanggan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pelanggan)
}

// DELETE /pelanggan/:id
func DeletePelanggan(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var prof models.Profil
	if err := db.Where("user_id = ?", id).First(&prof).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data profil pelanggan tidak ditemukan"})
		return
	}

	if err := db.Delete(&prof).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghapus data"})
		return
	}

	if err := db.Delete(&models.User{}, prof.UserID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghapus data user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User dengan role pelanggan terkait berhasil dihapus"})
}
