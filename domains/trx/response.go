package trx

import (
	"github.com/devanadindraa/Evermos-Backend/domains/address"
	"github.com/devanadindraa/Evermos-Backend/domains/product"
	"github.com/devanadindraa/Evermos-Backend/domains/shop"
)

type TrxRes struct {
	ID          int                 `json:"id"`
	HargaTotal  int                 `json:"harga_total"`
	KodeInvoice string              `json:"kode_invoice"`
	MethodBayar string              `json:"method_bayar"`
	AlamatKirim *address.AddressRes `json:"alamat_kirim"`
	DetailTrx   []DetailTrxRes      `json:"detail_trx"`
}

type DetailTrxRes struct {
	Product    product.ProductRes `json:"product"`
	Toko       *shop.ShopRes      `json:"toko"`
	Kuantitas  int                `json:"kuantitas"`
	HargaTotal int                `json:"harga_total"`
}

type PaginatedTrxRes struct {
	Data  []TrxRes `json:"data"`
	Page  int      `json:"page"`
	Limit int      `json:"limit"`
}
