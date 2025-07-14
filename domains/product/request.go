package product

import (
	"mime/multipart"

	"github.com/devanadindraa/Evermos-Backend/utils/constants"
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

type UpdateProductReq struct {
	NamaProduk    *string                  `form:"nama_produk"`
	Slug          *string                  `form:"slug"`
	IdCategory    *int                     `form:"category_id"`
	HargaReseller *string                  `form:"harga_reseller"`
	HargaKonsumen *string                  `form:"harga_konsumen"`
	Stok          *int                     `form:"stok"`
	Deskripsi     *string                  `form:"deskripsi"`
	Photos        *[]*multipart.FileHeader `form:"photos"`
}

type GetProductReq struct {
	*constants.FilterReq
	CategoryID *uint `query:"category_id"`
	TokoID     *uint `query:"toko_id"`
	MinHarga   *int  `query:"min_harga"`
	MaxHarga   *int  `query:"max_harga"`
}
