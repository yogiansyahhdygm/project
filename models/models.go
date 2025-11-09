package models

import (
	"time"
)

// User model matches users table
type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"size:50;unique;not null" json:"username"`
	Password string `gorm:"size:255;not null" json:"password,omitempty"`
	Role     string `gorm:"size:20;default:pelanggan" json:"role"`
}

// Kategori
type Kategori struct {
	ID      uint     `gorm:"primaryKey" json:"id"`
	Nama    string   `gorm:"size:100;not null" json:"nama"`
	Produks []Produk `gorm:"foreignKey:KategoriID" json:"produks,omitempty"`
}

// Produk
type Produk struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Nama       string    `gorm:"size:100;not null" json:"nama"`
	Harga      float64   `gorm:"type:decimal(10,2);not null" json:"harga"`
	Stok       int       `gorm:"not null" json:"stok"`
	KategoriID *uint     `json:"kategori_id"`
	Kategori   *Kategori `json:"kategori,omitempty"`
}

// Profil
type Profil struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Nama   string `gorm:"size:100" json:"nama"`
	Email  string `gorm:"size:100" json:"email"`
	Alamat string `gorm:"type:text" json:"alamat"`
	UserID *uint  `json:"user_id"`
	User   *User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"user,omitempty"`
}

// Pesanan
type Pesanan struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ProfilID   *uint     `json:"profil_id"`
	Profil     *Profil   `json:"profil,omitempty"`
	ProdukID   *uint     `json:"produk_id"`
	Produk     *Produk   `json:"produk,omitempty"`
	Jumlah     int       `gorm:"not null" json:"jumlah"`
	TotalHarga float64   `gorm:"type:decimal(12,2);not null" json:"total_harga"`
	Tanggal    time.Time `gorm:"autoCreateTime" json:"tanggal"`
	Status     string    `gorm:"type:varchar(20);default:'pending'" json:"status"`
}

func (Produk) TableName() string {
	return "produk"
}

func (Profil) TableName() string {
	return "profil"
}

func (Pesanan) TableName() string {
	return "pesanan"
}

func (Kategori) TableName() string {
	return "kategori"
}
