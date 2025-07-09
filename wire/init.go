package wireinject

import (
	"io"

	"github.com/devanadindraa/Evermos-Backend/routes"
	"github.com/devanadindraa/Evermos-Backend/utils/config"
	"github.com/devanadindraa/Evermos-Backend/utils/logger"
	"github.com/sirupsen/logrus"
)

func getLevel(level string) logrus.Level {
	switch level {
	case "TRACE":
		return logrus.TraceLevel
	case "DEBUG":
		return logrus.DebugLevel
	case "INFO":
		return logrus.InfoLevel
	case "WARN":
		return logrus.WarnLevel
	case "ERROR":
		return logrus.ErrorLevel
	case "FATAL":
		return logrus.FatalLevel
	case "PANIC":
		return logrus.PanicLevel
	}
	panic("invalid logger level")
}

func InitializeDependency(conf *config.Config) (*routes.Dependency, error) {
	logger.Setdata(conf.Environment.String(), conf.Version)
	logrus.SetLevel(getLevel(conf.Logger.Level))
	logrus.SetFormatter(&logrus.JSONFormatter{})
	if conf.Environment == config.TEST_ENVIRONMENT {
		logrus.SetOutput(io.Discard)
	}
	return initializeDependency(conf)
}
