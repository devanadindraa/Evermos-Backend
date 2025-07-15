package trx

import "time"

type Trx struct {
	ID               uint      `gorm:"primaryKey"`
	IdUser           uint      `gorm:"not null"`
	AlamatPengiriman uint      `gorm:"not null"`
	HargaTotal       int       `json:"harga_total"`
	KodeInvoice      string    `json:"kode_invoice"`
	MethodBayar      string    `json:"method_bayar"`
	CreatedAtDate    time.Time `gorm:"autoCreateTime"`
	UpdatedAtDate    time.Time `gorm:"autoUpdateTime"`
}

type LogProduk struct {
	ID            uint      `gorm:"primaryKey"`
	IdProduk      uint      `gorm:"not null"`
	NamaProduk    string    `json:"nama_produk"`
	Slug          string    `json:"slug"`
	HargaReseller string    `json:"harga_reseller"`
	HargaKonsumen string    `json:"harga_konsumen"`
	Deskripsi     string    `json:"deskripsi"`
	CreatedAtDate time.Time `gorm:"autoCreateTime"`
	UpdatedAtDate time.Time `gorm:"autoUpdateTime"`
	IdToko        uint      `gorm:"not null"`
	IdCategory    uint      `gorm:"not null"`
}

type DetailTrx struct {
	ID            uint      `gorm:"primaryKey"`
	IdTrx         uint      `gorm:"not null"`
	IdLogProduk   uint      `gorm:"not null"`
	IdToko        uint      `gorm:"not null"`
	Kuantitas     int       `json:"kuantitas"`
	HargaTotal    int       `json:"harga_total"`
	CreatedAtDate time.Time `gorm:"autoCreateTime"`
	UpdatedAtDate time.Time `gorm:"autoUpdateTime"`
}

func (Trx) TableName() string {
	return "trx"
}

func (LogProduk) TableName() string {
	return "log_produk"
}
func (DetailTrx) TableName() string {
	return "detail_trx"
}
