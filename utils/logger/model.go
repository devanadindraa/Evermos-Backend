package logger

import (
	"time"

	"github.com/devanadindraa/Evermos-Backend/utils/constants"
)

const PACKAGE_NAME = "github.com/devanadindraa/Evermos-Backend"

type LogPayload struct {
	Method         string
	Path           string
	StatusCode     int
	Took           time.Duration
	RequestPayload *constants.RequestPayload
}

func Setdata(env, ver string) {
	environment = env
	version = ver
}

var (
	environment = "unknown"
	version     = "unknown"
)
