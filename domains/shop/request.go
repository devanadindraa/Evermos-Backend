package shop

import "mime/multipart"

type UpdateShopReq struct {
	NamaToko *string               `form:"nama_toko"`
	UrlFoto  *multipart.FileHeader `form:"photo"`
}
