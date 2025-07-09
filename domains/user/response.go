package user

import (
	"time"
)

type LoginRes struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

type VerifyTokenRes struct {
	TokenVerified bool `json:"tokenVerified"`
}

type LogoutRes struct {
	LoggedOut bool `json:"loggedOut"`
}
