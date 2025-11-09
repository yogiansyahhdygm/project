package controllers

import (
	"fmt"
	"net/http"

	"toko/models"
	"toko/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreatePesanan(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	claims := c.MustGet("claims").(*utils.Claims)
	userID := uint(claims.UserID)

	var req struct {
		ProdukID *uint `json:"produk_id" binding:"required"`
		Jumlah   int   `json:"jumlah" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		// Pastikan profil ada
		var profil models.Profil
		if err := tx.First(&profil, userID).Error; err != nil {
			profil = models.Profil{ID: userID}
			if err := tx.Create(&profil).Error; err != nil {
				return err
			}
		}

		// Cek produk dan stok
		var produk models.Produk
		if err := tx.First(&produk, *req.ProdukID).Error; err != nil {
			return fmt.Errorf("produk tidak ditemukan")
		}
		if produk.Stok < req.Jumlah {
			return fmt.Errorf("stok tidak mencukupi")
		}

		// Kurangi stok
		produk.Stok -= req.Jumlah
		if err := tx.Save(&produk).Error; err != nil {
			return err
		}

		total := float64(req.Jumlah) * produk.Harga
		pesanan := models.Pesanan{
			ProfilID:   &profil.ID,
			ProdukID:   req.ProdukID,
			Jumlah:     req.Jumlah,
			TotalHarga: total,
			Status:     "pending",
		}
		if err := tx.Create(&pesanan).Error; err != nil {
			return err
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Pesanan berhasil dibuat",
			"data":    pesanan,
		})
		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// Pesanan saya (pelanggan)
func GetMyPesanan(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	claims := c.MustGet("claims").(*utils.Claims)
	userID := uint(claims.UserID)

	var pesanan []models.Pesanan
	if err := db.
		Where("profil_id = ?", userID).
		Preload("Produk").
		Order("tanggal DESC").
		Find(&pesanan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": pesanan})
}

// all pesanan - admin
func GetAllPesanan(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var pesanan []models.Pesanan
	if err := db.
		Preload("Profil").
		Preload("Produk").
		Order("tanggal DESC").
		Find(&pesanan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": pesanan})
}

// Admin edit status pemesanan
func UpdateStatusPesanan(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		var pesanan models.Pesanan
		if err := tx.Preload("Produk").First(&pesanan, id).Error; err != nil {
			return fmt.Errorf("pesanan tidak ditemukan")
		}

		// Update status & kelola stok
		switch req.Status {
		case "dibatalkan":
			// Kembalikan stok produk
			pesanan.Produk.Stok += pesanan.Jumlah
			if err := tx.Save(pesanan.Produk).Error; err != nil {
				return err
			}
			pesanan.Status = "dibatalkan"

		case "selesai":
			// Tidak ubah stok lagi, hanya ubah status
			pesanan.Status = "selesai"

		case "pending":
			pesanan.Status = "pending"

		default:
			return fmt.Errorf("status tidak valid (gunakan pending, selesai, atau dibatalkan)")
		}

		if err := tx.Save(&pesanan).Error; err != nil {
			return err
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Status pesanan diperbarui ke '%s'", req.Status),
			"data":    pesanan,
		})
		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func DeleteMyPesanan(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	claims := c.MustGet("claims").(*utils.Claims)
	userID := uint(claims.UserID)

	err := db.Transaction(func(tx *gorm.DB) error {
		var pesanan models.Pesanan

		// Cari pesanan berdasarkan id & profil_id (harus milik user)
		if err := tx.Preload("Produk").
			Where("id = ? AND profil_id = ?", id, userID).
			First(&pesanan).Error; err != nil {
			return fmt.Errorf("pesanan tidak ditemukan atau bukan milik kamu")
		}

		// Jika pending → kembalikan stok produk
		if pesanan.Status == "pending" {
			pesanan.Produk.Stok += pesanan.Jumlah
			if err := tx.Save(pesanan.Produk).Error; err != nil {
				return err
			}
		}

		// Hapus pesanan
		if err := tx.Delete(&pesanan).Error; err != nil {
			return err
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Pesanan kamu berhasil dihapus",
			"data":    pesanan,
		})
		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func DeletePesananAdmin(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	err := db.Transaction(func(tx *gorm.DB) error {
		var pesanan models.Pesanan

		if err := tx.Preload("Produk").First(&pesanan, id).Error; err != nil {
			return fmt.Errorf("pesanan tidak ditemukan")
		}

		// Jika pending → kembalikan stok produk
		if pesanan.Status == "pending" {
			pesanan.Produk.Stok += pesanan.Jumlah
			if err := tx.Save(pesanan.Produk).Error; err != nil {
				return err
			}
		}

		// Hapus pesanan
		if err := tx.Delete(&pesanan).Error; err != nil {
			return err
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Pesanan berhasil dihapus oleh admin",
			"data":    pesanan,
		})
		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
