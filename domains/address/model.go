package address

import "time"

type Address struct {
	ID            uint      `gorm:"primaryKey"`
	IdUser        uint      `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	NoTelp        string    `json:"no_telp"`
	JudulAlamat   string    `json:"judul_alamat"`
	NamaPenerima  string    `json:"nama_penerima"`
	DetailAlamat  string    `json:"detail_alamat"`
	CreatedAtDate time.Time `gorm:"autoCreateTime"`
	UpdatedAtDate time.Time `gorm:"autoUpdateTime"`
}

func (Address) TableName() string {
	return "alamat"
}
