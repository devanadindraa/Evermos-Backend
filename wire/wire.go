//go:build wireinject
// +build wireinject

package wireinject

import (
	"github.com/devanadindraa/Evermos-Backend/database"
	"github.com/devanadindraa/Evermos-Backend/domains/address"
	"github.com/devanadindraa/Evermos-Backend/domains/category"
	"github.com/devanadindraa/Evermos-Backend/domains/product"
	"github.com/devanadindraa/Evermos-Backend/domains/provcity"
	"github.com/devanadindraa/Evermos-Backend/domains/shop"
	"github.com/devanadindraa/Evermos-Backend/domains/trx"
	"github.com/devanadindraa/Evermos-Backend/domains/user"
	"github.com/devanadindraa/Evermos-Backend/middlewares"
	"github.com/devanadindraa/Evermos-Backend/routes"
	"github.com/devanadindraa/Evermos-Backend/utils/config"
	"github.com/go-playground/validator/v10"
	_ "github.com/google/subcommands"
	"github.com/google/wire"
)

var userSet = wire.NewSet(
	user.NewService,
	user.NewHandler,
)

var provcitySet = wire.NewSet(
	provcity.NewEmsiaClient,
	provcity.NewProvcityHandler,
)

var categorySet = wire.NewSet(
	category.NewService,
	category.NewHandler,
)

var shopSet = wire.NewSet(
	shop.NewService,
	shop.NewHandler,
)

var addressSet = wire.NewSet(
	address.NewService,
	address.NewHandler,
)

var productSet = wire.NewSet(
	product.NewService,
	product.NewHandler,
)

var trxSet = wire.NewSet(
	trx.NewService,
	trx.NewHandler,
)

func NewValidator() *validator.Validate {
	return validator.New()
}

func initializeDependency(config *config.Config) (*routes.Dependency, error) {

	wire.Build(
		database.NewDB,
		middlewares.NewMiddlewares,
		NewValidator,
		routes.NewDependency,
		userSet,
		provcitySet,
		categorySet,
		shopSet,
		addressSet,
		productSet,
		trxSet,
	)

	return nil, nil
}
