package routes

import (
	"github.com/devanadindraa/Evermos-Backend/middlewares"
	"github.com/devanadindraa/Evermos-Backend/utils/config"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/cors"
	"gorm.io/gorm"
)

func NewDependency(
	conf *config.Config,
	mw middlewares.Middlewares,
	db *gorm.DB,
) *Dependency {

	router := fiber.New()
	router.Use(func(c *fiber.Ctx) error {
		if c.Method() != fiber.MethodGet &&
			c.Method() != fiber.MethodPost &&
			c.Method() != fiber.MethodPut &&
			c.Method() != fiber.MethodDelete {
			return c.Status(fiber.StatusMethodNotAllowed).SendString("Method Not Allowed")
		}
		return c.Next()
	})

	// middleware
	{
		router.Use(cors.Default())
		router.Use(mw.AddRequestId)
		router.Use(mw.Logging)
		router.Use(mw.RateLimiter)
		router.Use(mw.Recover)
	}

	return &Dependency{
		handler: router,
		db:      db,
	}
}
