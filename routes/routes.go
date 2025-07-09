package routes

import (
	"github.com/devanadindraa/Evermos-Backend/domains/user"
	"github.com/devanadindraa/Evermos-Backend/middlewares"
	"github.com/devanadindraa/Evermos-Backend/utils/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"gorm.io/gorm"
)

func NewDependency(
	conf *config.Config,
	mw middlewares.Middlewares,
	db *gorm.DB,
	userHandler user.Handler,
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
		router.Use(cors.New())
		router.Use(mw.AddRequestId)
		router.Use(mw.Logging)
		router.Use(mw.RateLimiter)
		router.Use(mw.Recover)
	}

	// domain user
	auth := router.Group("/auth")
	{
		auth.Post("/login", mw.BasicAuth, userHandler.Login)
		auth.Get("/verify-token", mw.JWT, userHandler.VerifyToken)
		auth.Post("/logout", mw.JWT, userHandler.Logout)
		auth.Post("/register", mw.BasicAuth, userHandler.Register)
	}

	return &Dependency{
		handler: router,
		db:      db,
	}
}
