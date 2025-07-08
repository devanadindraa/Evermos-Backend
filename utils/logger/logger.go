package logger

import (
	"context"
	"fmt"
	"time"

	apierror "github.com/devanadindraa/Evermos-Backend/utils/api-error"
	contextUtil "github.com/devanadindraa/Evermos-Backend/utils/context"
	"github.com/sirupsen/logrus"
)

func Trace(ctx context.Context, format string, values ...interface{}) {
	logrus.
		WithContext(ctx).
		WithFields(logrus.Fields{
			"requestId":   contextUtil.GetRequestId(ctx),
			"callers":     getCaller(),
			"environment": environment,
			"version":     version,
		}).Trace(fmt.Sprintf(format, values...))
}

func Info(ctx context.Context, format string, values ...interface{}) {
	logrus.
		WithContext(ctx).
		WithFields(logrus.Fields{
			"requestId":   contextUtil.GetRequestId(ctx),
			"callers":     getCaller(),
			"environment": environment,
			"version":     version,
		}).Info(fmt.Sprintf(format, values...))
}

func Error(ctx context.Context, format string, values ...interface{}) {
	logrus.
		WithContext(ctx).
		WithFields(logrus.Fields{
			"requestId":   contextUtil.GetRequestId(ctx),
			"callers":     getCaller(),
			"environment": environment,
			"version":     version,
		}).Error(fmt.Sprintf(format, values...))
}

func Warn(ctx context.Context, format string, values ...interface{}) {
	logrus.
		WithContext(ctx).
		WithFields(logrus.Fields{
			"requestId":   contextUtil.GetRequestId(ctx),
			"callers":     getCaller(),
			"environment": environment,
			"version":     version,
		}).Warn(fmt.Sprintf(format, values...))
}

func TraceErr(ctx context.Context, err error) {
	err, callers := getTracerrCallers(err)
	fields := logrus.Fields{
		"requestId":   contextUtil.GetRequestId(ctx),
		"callers":     callers,
		"environment": environment,
		"version":     version,
	}
	logrus.
		WithContext(ctx).
		WithFields(fields).
		Trace(err.Error())
}

func Log(ctx context.Context, payload LogPayload, err error) {

	// use manual convert cause regular func return int not float
	tookInMillisecond := float64(payload.Took) / float64(time.Millisecond)

	fields := logrus.Fields{
		"requestId":         contextUtil.GetRequestId(ctx),
		"callers":           make([]string, 0),
		"statusCode":        payload.StatusCode,
		"method":            payload.Method,
		"path":              payload.Path,
		"tookInMillisecond": tookInMillisecond,
		"payload":           payload.RequestPayload,
		"environment":       environment,
		"version":           version,
	}
	if err == nil {
		logrus.
			WithContext(ctx).
			WithFields(fields).
			Info("")
		return
	}

	err, callers := getTracerrCallers(err)

	fields["callers"] = callers
	message := err.Error()
	apiErrors, ok := err.(apierror.ApiErrors)
	if !ok {
		logrus.
			WithContext(ctx).
			WithFields(fields).
			Error(message)
		return
	}

	fields["statusCode"] = apiErrors.Code
	switch apiErrors.Level {
	case "ERROR":
		logrus.
			WithContext(ctx).
			WithFields(fields).
			Error(message)
	case "WARN":
		logrus.
			WithContext(ctx).
			WithFields(fields).
			Warn(message)
	}

}
