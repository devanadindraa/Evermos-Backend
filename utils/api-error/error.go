package apierror

import (
	"net/http"
)

func Unauthorized() error {
	return NewWarn(http.StatusUnauthorized, "Unauthorized!")
}

func FailedToConvertUpdatedAt() error {
	return NewError(http.StatusInternalServerError, "Failed convert updated at to time")
}

func FailedToConvertCreatedAt() error {
	return NewError(http.StatusInternalServerError, "Failed convert created at to time")
}

func FileNotFound() error {
	return NewWarn(http.StatusNotFound, "File not found!")
}

func InvalidFileId() error {
	return NewWarn(http.StatusBadRequest, "fileId must be UUID!")
}
