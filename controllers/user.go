package controllers

import (
	"fmt"
	"net/http"

	"toko/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// GET all /user
func GetAllUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var us []models.User
	if err := db.Find(&us).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, us)
}

// DELETE /User/:id
func DeleteUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data User tidak ditemukan"})
		return
	}

	var pr models.Profil

	if err := db.Where("user_id = ?", user.ID).First(&pr).Error; err == nil {
		// hapus pesanan yang terkait profil
		if err := db.Where("profil_id = ?", pr.ID).Delete(&models.Pesanan{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Errorf("gagal menghapus pesanan: %v", err),
			})
			return
		}
		// hapus profil
		if err := db.Delete(&pr).Error; err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Errorf("gagal menghapus profil: %v", err),
			})
			return
		}
	}
	//hapus user
	if err := db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghapus data user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User dengan role pelanggan terkait berhasil dihapus"})
}

// mengganti hak akses
func EditdataUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var us models.User
	if err := db.First(&us, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan!"})
		return
	}

	if err := c.ShouldBindJSON(&us); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal payload"})
		return
	}

	updateData := map[string]interface{}{}
	if us.Username != "" {
		updateData["username"] = us.Username
	}
	if us.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(us.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengenkripsi password"})
			return
		}
		updateData["password"] = hashed
	}
	if us.Role != "" {
		updateData["role"] = us.Role
	}

	if err := db.Model(&us).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data User berhasil diperbarui",
		"User":    us,
	})
}
