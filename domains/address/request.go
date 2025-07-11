package address

type AddressReq struct {
	NoTelp       *string `json:"no_telp"`
	JudulAlamat  string  `json:"judul_alamat" validated:"required"`
	NamaPenerima string  `json:"nama_penerima" validated:"required"`
	DetailAlamat string  `json:"detail_alamat" validated:"required"`
}

type UpdateAddressReq struct {
	NoTelp       *string `json:"no_telp"`
	JudulAlamat  *string `json:"judul_alamat"`
	NamaPenerima *string `json:"nama_penerima"`
	DetailAlamat *string `json:"detail_alamat"`
}
