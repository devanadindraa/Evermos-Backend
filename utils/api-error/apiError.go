package apierror

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/ztrue/tracerr"
)

type ApiErrors struct {
	Code     int
	Level    errorLevel
	Messages []string
}

func (apiErrors ApiErrors) Error() string {
	errorStr := ""
	for index, message := range apiErrors.Messages {
		if index != 0 {
			errorStr += " - "
		}
		errorStr += message
	}
	return errorStr
}

func Warn(code int, err error) error {
	return NewWarn(code, err.Error())
}

func NewWarn(code int, errMessages ...string) error {
	return tracerr.Wrap(ApiErrors{Code: code, Level: WARN_LEVEL, Messages: errMessages})
}

func Error(code int, err error) error {
	return NewError(code, err.Error())
}

func NewError(code int, errMessages ...string) error {
	return tracerr.Wrap(ApiErrors{Code: code, Level: ERROR_LEVEL, Messages: errMessages})
}

func FromErr(err error) error {
	if err == nil {
		return nil
	}
	var (
		apiErrors ApiErrors
		ve        validator.ValidationErrors
	)
	if errors.As(tracerr.Unwrap(err), &apiErrors) {
		return err
	} else if errors.As(tracerr.Unwrap(err), &ve) {
		messages := make([]string, 0)
		for _, fe := range ve {
			messages = append(messages, fmt.Sprintf("%s: %s", fe.Field(), msgForTag(fe.Type().Kind(), fe.Tag(), fe.Param())))
		}
		return NewWarn(http.StatusBadRequest, messages...)
	} else {
		return Error(http.StatusInternalServerError, err)
	}
}

func msgForTag(kind reflect.Kind, tag, param string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "number":
		return "This field must be number"
	case "gt":
		switch kind {
		case reflect.Array, reflect.Slice:
			return fmt.Sprintf("Array length must be greater than %s", param)
		default:
			return fmt.Sprintf("This field must be greater than %s", param)
		}
	}
	return fmt.Sprintf("Failed on %s tag", tag)
}

func RBError(rbErr error, err error) ApiErrors {
	if err == nil {
		return ApiErrors{Code: http.StatusInternalServerError, Level: ERROR_LEVEL, Messages: []string{rbErr.Error()}}
	}
	var apiErrors ApiErrors
	if errors.As(err, &apiErrors) {
		return ApiErrors{Code: http.StatusInternalServerError, Level: ERROR_LEVEL, Messages: append(apiErrors.Messages, rbErr.Error())}
	} else {
		return ApiErrors{Code: http.StatusInternalServerError, Level: ERROR_LEVEL, Messages: []string{err.Error(), rbErr.Error()}}
	}
}

func GetApiErrors(err error) ApiErrors {
	if err == nil {
		return ApiErrors{
			Code:     http.StatusInternalServerError,
			Level:    ERROR_LEVEL,
			Messages: []string{"empty error"},
		}
	}

	var apiErrors ApiErrors
	if errors.As(tracerr.Unwrap(err), &apiErrors) {
		return apiErrors
	}

	return ApiErrors{
		Code:     http.StatusInternalServerError,
		Level:    ERROR_LEVEL,
		Messages: []string{err.Error()},
	}
}

type ErrGroup struct {
	mu      sync.Mutex
	errsStr []string
}

func (e *ErrGroup) Append(errStr string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.errsStr = append(e.errsStr, errStr)
}

func (e *ErrGroup) GetErr() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if len(e.errsStr) == 0 {
		return nil
	}

	apiErrors := ApiErrors{
		Code:     http.StatusInternalServerError,
		Level:    ERROR_LEVEL,
		Messages: e.errsStr,
	}
	return tracerr.Wrap(apiErrors)
}
