package constants

import (
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTClaims struct {
	ID      int    `json:"userID"`
	IsAdmin bool   `json:"isAdmin"`
	NoTelp  string `json:"no_telp"`
	jwt.RegisteredClaims
}

type Token struct {
	Token  string
	Claims JWTClaims
}

type FilterReq struct {
	Limit          int64 `validate:"gte=1"`
	Page           int64 `validate:"gte=1"`
	OrderBy        string
	Keyword        string
	SortOrder      string `validate:"oneof=asc desc"`
	StartCreatedAt *time.Time
	EndCreatedAt   *time.Time
	StartUpdatedAt *time.Time
	EndUpdatedAt   *time.Time
}

type MetaData struct {
	Page      int64 `json:"page"`
	TotalPage int64 `json:"totalPage"`
	TotalData int64 `json:"totalData"`
}

type Pagination[T any] struct {
	Data       T        `json:"data"`
	Pagination MetaData `json:"pagination"`
}

type RequestPayload struct {
	Body        map[string]any      `json:"body" bson:"body"`
	QueryParams url.Values          `json:"queryParams" bson:"queryParams"`
	Headers     map[string][]string `json:"headers" bson:"headers"`
}
