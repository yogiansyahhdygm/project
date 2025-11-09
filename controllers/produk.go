package controllers

import (
	"net/http"

	"toko/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateProduk(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var req models.Produk
	if err := c.ShouldBindJSON(&req); err != nil || req.Nama == "" || req.Harga == 0 || req.Stok == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama, Harga, dan Stok tidak boleh kosong!"})
		return
	}

	p := models.Produk{
		Nama:       req.Nama,
		Harga:      req.Harga,
		Stok:       req.Stok,
		KategoriID: req.KategoriID,
	}
	if err := db.Create(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func ListProduk(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var prods []models.Produk
	db.Preload("Kategori").Find(&prods)
	c.JSON(http.StatusOK, prods)
}

func GetProdukByID(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var p models.Produk
	if err := db.First(&p, id).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Produk Tidak ditemukan!!"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func UpdateProduk(c *gin.Context) {
	id := c.Param("id")
	db := c.MustGet("db").(*gorm.DB)

	var p models.Produk
	if err := db.First(&p, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan!"})
		return
	}

	var input struct {
		Nama       string `json:"nama"`
		Harga      int    `json:"harga"`
		Stok       int    `json:"stok"`
		KategoriID int    `json:"kategori_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	updateData := map[string]interface{}{}
	if input.Nama != "" {
		updateData["nama"] = input.Nama
	}
	if input.Harga != 0 {
		updateData["harga"] = input.Harga
	}
	if input.Stok != 0 {
		updateData["stok"] = input.Stok
	}
	if input.KategoriID != 0 {
		var kt models.Kategori
		if err := db.Where("id = ?", input.KategoriID).First(&kt).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data kategori tidak ditemukan"})
			return
		}
		updateData["kategori_id"] = input.KategoriID
	}

	if err := db.Model(&p).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "data profil berhasil diperbarui",
		"Produk":  p,
	})
}

func DeleteProduk(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var p models.Produk
	if err := db.Where("id = ?", id).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data Produk tidak ditemukan"})
		return
	}

	if err := db.Delete(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghapus data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Produk terkait berhasil dihapus"})
}
