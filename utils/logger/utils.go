package logger

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/ztrue/tracerr"
)

func getCaller() []string {
	pc, _, line, ok := runtime.Caller(2)
	if !ok {
		line = 0
	}

	funcName := runtime.FuncForPC(pc).Name()

	return []string{
		fmt.Sprintf(
			"%s:%d",
			strings.TrimPrefix(
				funcName,
				PACKAGE_NAME,
			),
			line,
		),
	}
}

func getTracerrCallers(err error) (originalErr error, callers []string) {
	tracerErr, ok := err.(tracerr.Error)
	if !ok {
		return err, make([]string, 0)
	}

	for _, frame := range tracerErr.StackTrace() {

		funcSplit := strings.Split(frame.Func, "github.com/devanadindraa/Evermos-Backend")
		if len(funcSplit) == 1 {
			continue
		}

		funcSplit = strings.Split(funcSplit[1], "/utils/api-error")
		if len(funcSplit) == 2 {
			continue
		}

		callers = append(callers, fmt.Sprintf("%s:%d", funcSplit[0], frame.Line))
	}

	return tracerErr.Unwrap(), callers
}
