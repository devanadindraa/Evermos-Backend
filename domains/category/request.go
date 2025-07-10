package category

type CategoryReq struct {
	NamaCategory string `json:"nama_category" validate:"required"`
}
