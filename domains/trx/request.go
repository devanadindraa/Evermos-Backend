package trx

type TrxReq struct {
	MethodBayar string         `json:"method_bayar"`
	AlamatKirim int            `json:"alamat_kirim"`
	DetailTrx   []DetailTrxReq `json:"detail_trx"`
}

type DetailTrxReq struct {
	ProdukId  int `json:"product_id"`
	Kuantitas int `json:"kuantitas"`
}
