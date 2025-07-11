package product

import (
	"time"
)

type Product struct {
	ID            uint    `gorm:"primaryKey"`
	IdToko        uint    `gorm:"not null"`
	NamaProduk    string  `json:"nama_produk"`
	IdCategory    uint    `gorm:"not null"`
	Slug          string  `json:"slug"`
	HargaReseller string  `json:"harga_reseller"`
	HargaKonsumen string  `json:"harga_konsumen"`
	Stok          int     `json:"stok"`
	Deskripsi     string  `json:"deskripsi"`
	Photos        []Photo `gorm:"foreignKey:IdProduk;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"photos"`

	CreatedAtDate time.Time `gorm:"autoCreateTime"`
	UpdatedAtDate time.Time `gorm:"autoUpdateTime"`
}

type Photo struct {
	ID       uint   `gorm:"primaryKey"`
	IdProduk uint   `gorm:"not null"`
	Url      string `json:"url"`

	CreatedAtDate time.Time `gorm:"autoCreateTime"`
	UpdatedAtDate time.Time `gorm:"autoUpdateTime"`
}

func (Product) TableName() string {
	return "produk"
}

func (Photo) TableName() string {
	return "foto_produk"
}
