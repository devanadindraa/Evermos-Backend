package product

import (
	"github.com/devanadindraa/Evermos-Backend/domains/category"
	"github.com/devanadindraa/Evermos-Backend/domains/shop"
)

type ProductRes struct {
	ID            int                   `json:"id"`
	NamaProduk    *string               `json:"nama_produk,omitempty"`
	Slug          *string               `json:"slug,omitempty"`
	IdCategory    *uint                 `json:"category_id,omitempty"`
	HargaReseller *string               `json:"harga_reseller,omitempty"`
	HargaKonsumen *string               `json:"harga_konsumen,omitempty"`
	Stok          *int                  `json:"stok,omitempty"`
	Deskripsi     *string               `json:"deskripsi,omitempty"`
	Shop          *shop.ShopRes         `json:"shop,omitempty"`
	Category      *category.CategoryRes `json:"category,omitempty"`
	Photos        []string              `json:"photos,omitempty"`
}
