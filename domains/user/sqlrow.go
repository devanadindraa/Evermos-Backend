package user

import (
	"time"
)

type InvalidToken struct {
	Token   string
	Expires time.Time
}

func (InvalidToken) TableName() string {
	return "invalid_token"
}

type User struct {
	ID            uint      `gorm:"primaryKey;autoIncrement"`
	Nama          string    `json:"nama"`
	KataSandi     string    `json:"kata_sandi"`
	Notelp        string    `json:"notelp" gorm:"unique"`
	TanggalLahir  time.Time `json:"tanggal_Lahir" gorm:"type:date"`
	Pekerjaan     string    `json:"pekerjaan"`
	Email         string    `json:"email"`
	IdProvinsi    string    `json:"id_provinsi"`
	IdKota        string    `json:"id_kota"`
	IsAdmin       bool      `json:"isAdmin" gorm:"column:isAdmin;default:false"`
	CreatedAtDate time.Time `gorm:"autoCreateTime"`
	UpdatedAtDate time.Time `gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "user"
}
