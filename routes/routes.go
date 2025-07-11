package routes

import (
	"github.com/devanadindraa/Evermos-Backend/domains/address"
	"github.com/devanadindraa/Evermos-Backend/domains/category"
	"github.com/devanadindraa/Evermos-Backend/domains/provcity"
	"github.com/devanadindraa/Evermos-Backend/domains/shop"
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
	provcityHandler provcity.Handler,
	categoryHandler category.Handler,
	shopHandler shop.Handler,
	addressHandler address.Handler,
) *Dependency {

	app := fiber.New()
	router := app.Group("/api/v1")
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

	// domain auth
	auth := router.Group("/auth")
	{
		auth.Post("/login", mw.BasicAuth, userHandler.Login)
		auth.Get("/verify-token", mw.JWT(false), userHandler.VerifyToken)
		auth.Post("/logout", mw.JWT(false), userHandler.Logout)
		auth.Post("/register", mw.BasicAuth, userHandler.Register)
	}

	// domain user
	user := router.Group("/user")
	{
		user.Put("", mw.JWT(false), userHandler.UpdateProfile)
		user.Get("", mw.JWT(false), userHandler.GetProfile)
		user.Post("/alamat", mw.JWT(false), addressHandler.AddAddress)
		user.Get("/alamat", mw.JWT(false), addressHandler.GetMyAddress)
	}

	// domain provcity
	provcity := router.Group("/provcity")
	{
		provcity.Get("/listprovincies", provcityHandler.GetProvinces)
		provcity.Get("/listcities/:prov_id", provcityHandler.GetCitys)
		provcity.Get("/detailprovince/:prov_id", provcityHandler.GetDetailProvince)
		provcity.Get("/detailcity/:city_id", provcityHandler.GetDetailCity)
	}

	// domain category
	category := router.Group("/category")
	{
		category.Post("", mw.JWT(true), categoryHandler.AddCategory)
		category.Get("", mw.JWT(true), categoryHandler.GetAllCategory)
		category.Get("/:id", mw.JWT(true), categoryHandler.GetCategoryByID)
		category.Delete("/:id", mw.JWT(true), categoryHandler.DeleteCategory)
		category.Put("/:id", mw.JWT(true), categoryHandler.UpdateCategory)
	}

	// domain toko
	shop := router.Group("/toko")
	{
		shop.Get("/my", mw.JWT(false), shopHandler.GetMyShop)
		shop.Get("/:id_toko", mw.JWT(false), shopHandler.GetShopByID)
		shop.Put("/:id_toko", mw.JWT(false), shopHandler.UpdateMyShop)
		shop.Get("/", mw.JWT(false), shopHandler.GetAllShop)
	}

	app.Use(func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Endpoint not found",
			"errors":  []string{"Please check the URL or HTTP method used"},
			"data":    nil,
		})
	})

	return &Dependency{
		handler: app,
		db:      db,
	}
}
