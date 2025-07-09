package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/devanadindraa/Evermos-Backend/utils/config"
	"github.com/devanadindraa/Evermos-Backend/utils/logger"
	wireinject "github.com/devanadindraa/Evermos-Backend/wire"

	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	conf := config.NewConfig()

	dependency, err := wireinject.InitializeDependency(conf)
	if err != nil {
		logger.Error(context.Background(), "failed to initialize dependency %v", err)
		panic(err)
	}
	defer dependency.Close()

	// Ambil instance fiber dari dependency
	app := dependency.GetHandler()

	// Tambahkan middleware CORS dari Fiber
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Authorization, Content-Type",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH",
	}))

	// Buat context yang mendeteksi interrupt signal (Ctrl+C)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Jalankan Fiber server di goroutine
	go func() {
		addr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
		if err := app.Listen(addr); err != nil {
			logger.Error(context.Background(), "Error starting Fiber server %v", err)
			panic(err)
		}
	}()

	// Tunggu sampai signal stop
	<-ctx.Done()
	stop()

	logger.Trace(context.Background(), "shutting down gracefully, please wait a moment ...")

	// Context timeout untuk shutdown
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Shutdown(); err != nil {
		logger.Error(timeoutCtx, "Server forced to shutdown %v", err)
	}

	logger.Trace(timeoutCtx, "Server gracefully shutdown")
}
