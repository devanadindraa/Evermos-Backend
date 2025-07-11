package shop

type ShopRes struct {
	ID       int    `json:"id"`
	NamaToko string `json:"nama_toko"`
	UrlFoto  string `json:"url_foto"`
	IdUser   *int   `json:"id_user,omitempty"`
}

type PaginatedShopRes struct {
	Page  int       `json:"page"`
	Limit int       `json:"limit"`
	Data  []ShopRes `json:"data"`
}
