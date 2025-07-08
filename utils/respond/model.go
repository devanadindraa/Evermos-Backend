package respond

type ApiModel[T any] struct {
	Status  bool     `json:"status"`
	Message string   `json:"message"`
	Errors  []string `json:"errors" example:""`
	Data    T        `json:"data"`
}

type DataParam struct {
	Code     int
	Filename string
	MimeType string
	Data     []byte
}
