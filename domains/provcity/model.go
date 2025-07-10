package provcity

type DefaultResponse[T any] struct {
	Meta Meta `json:"meta"`
	Data T    `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

type Province struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type City struct {
	ID         string `json:"id"`
	ProvinceId string `json:"province_id"`
	Name       string `json:"name"`
}
