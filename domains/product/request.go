package product

import (
	"mime/multipart"
)

type ProductReq struct {
	NamaProduk    string                  `form:"nama_produk" validate:"required"`
	Slug          *string                 `form:"slug"`
	IdCategory    uint                    `form:"category_id" validate:"required"`
	HargaReseller string                  `form:"harga_reseller" validate:"required"`
	HargaKonsumen string                  `form:"harga_konsumen" validate:"required"`
	Stok          int                     `form:"stok" validate:"required,min=0"`
	Deskripsi     string                  `form:"deskripsi" validate:"required"`
	Photos        []*multipart.FileHeader `form:"photos"`
}
