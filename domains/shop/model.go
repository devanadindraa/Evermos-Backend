package shop

import (
	"time"
)

type Toko struct {
	ID            uint `gorm:"primaryKey"`
	IdUser        uint `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	NamaToko      string
	UrlFoto       string
	CreatedAtDate time.Time
	UpdatedAtDate time.Time
}

func (Toko) TableName() string {
	return "toko"
}
