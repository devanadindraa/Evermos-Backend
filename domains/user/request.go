package user

import (
	"time"
)

type LoginReq struct {
	Notelp    string `json:"no_telp" validate:"required"`
	KataSandi string `json:"kata_sandi" validate:"required"`
}

type LogoutReq struct {
	Token   string
	Expires time.Time
}

type RegisterReq struct {
	Nama         string  `json:"nama" validate:"required"`
	KataSandi    string  `json:"kata_sandi" validate:"required"`
	NoTelp       string  `json:"no_telp" validate:"required"`
	TanggalLahir *string `json:"tanggal_Lahir"`
	JenisKelamin *string `json:"jenis_kelamin"`
	Tentang      *string `json:"tentang"`
	Pekerjaan    *string `json:"pekerjaan"`
	Email        string  `json:"email" validate:"required"`
	IdProvinsi   *string `json:"id_provinsi"`
	IdKota       *string `json:"id_kota"`
	IsAdmin      *bool   `json:"isAdmin,omitempty"`
}

type UpdateProfileReq struct {
	Nama         string  `json:"nama" validate:"required"`
	KataSandi    string  `json:"kata_sandi" validate:"required"`
	NoTelp       string  `json:"no_telp" validate:"required"`
	TanggalLahir *string `json:"tanggal_Lahir"`
	JenisKelamin *string `json:"jenis_kelamin"`
	Tentang      *string `json:"tentang"`
	Pekerjaan    *string `json:"pekerjaan"`
	Email        string  `json:"email" validate:"required"`
	IdProvinsi   *string `json:"id_provinsi"`
	IdKota       *string `json:"id_kota"`
	IsAdmin      *bool   `json:"isAdmin,omitempty"`
}
