package controllers

import (
	"net/http"

	"toko/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Allkategori(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var kt []models.Kategori
	if err := db.Find(&kt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kt)
}

func GetKategoriByID(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var ktr models.Kategori
	if err := db.First(&ktr, id).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Profil anda tidak ditemukan. Silahkan melakukan Register!"})
		return
	}
	c.JSON(http.StatusOK, ktr)
}

func CreateKategori(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var k models.Kategori
	if err := c.ShouldBindJSON(&k); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if k.Nama == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kategori tidak boleh kosong"})
		return
	}

	kategori := models.Kategori{
		Nama: k.Nama,
	}

	if err := db.Create(&kategori).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "kategori created"})
}

func UpdateKategori(c *gin.Context) {
	id := c.Param("id")
	db := c.MustGet("db").(*gorm.DB)

	var k models.Kategori
	if err := db.First(&k, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan!"})
		return
	}

	var input struct {
		Nama string `json:"nama"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	updateData := map[string]interface{}{}
	if input.Nama != "" {
		updateData["nama"] = input.Nama
	}

	if err := db.Model(&k).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "data profil berhasil diperbarui",
		"Profil":  k,
	})
}

func DeleteKategori(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var kt models.Kategori
	if err := db.Where("id = ?", id).First(&kt).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data kategori tidak ditemukan"})
		return
	}

	if err := db.Delete(&kt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghapus data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Kategori berhasil dihapus"})
}
