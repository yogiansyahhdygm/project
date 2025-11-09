package routes

import (
	"toko/controllers"
	"toko/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(r *gin.Engine, db *gorm.DB) {
	// inject db into context
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	//Selamat Datang
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome to Toko API")
	})

	// public auth
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// public product listing
	r.GET("/produk", controllers.ListProduk)
	r.GET("/produk/:id", controllers.GetProdukByID)

	// pelanggan routes (JWT)
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// profile endpoints (pelanggan)
		api.GET("/profil", controllers.GetProfil)
		api.PUT("/profil", controllers.UpdateProfil)

		// pesanan (buat, lihat, hapus)
		api.GET("/pesanan", controllers.GetMyPesanan)
		api.POST("/pesanan", controllers.CreatePesanan)
		api.DELETE("/pesanan/:id", controllers.DeleteMyPesanan)
	}

	// admin group (JWT + AdminOnly)
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminOnly())
	{
		// produk
		admin.GET("/produk", controllers.ListProduk)
		admin.GET("/produk/:id", controllers.GetProdukByID)
		admin.POST("/produk", controllers.CreateProduk)
		admin.PUT("/produk/:id", controllers.UpdateProduk)
		admin.DELETE("/produk/:id", controllers.DeleteProduk)

		//kategori
		admin.GET("/kategori", controllers.Allkategori)
		admin.GET("/kategori/:id", controllers.GetKategoriByID)
		admin.POST("/kategori", controllers.CreateKategori)
		admin.PUT("/kategori/:id", controllers.UpdateKategori)
		admin.DELETE("/kategori/:id", controllers.DeleteKategori)

		//pesanan (lihat dan update status)
		admin.GET("/pesanan", controllers.GetAllPesanan)
		admin.PUT("/pesanan/:id", controllers.UpdateStatusPesanan)
		admin.DELETE("/pesanan/:id", controllers.DeletePesananAdmin)

		//Pelanggan
		admin.GET("/pelanggan", controllers.GetAllPelanggan)
		admin.DELETE("/pelanggan/:id", controllers.DeletePelanggan)

		//user - lihat all user, ganti hak akses, delete user
		admin.GET("/user", controllers.GetAllUser)
		admin.DELETE("/user/:id", controllers.DeleteUser)
		admin.PUT("/user/:id", controllers.EditdataUser)

	}
}
