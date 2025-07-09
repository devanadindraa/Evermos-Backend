//go:build wireinject
// +build wireinject

package wireinject

import (
	"github.com/devanadindraa/Evermos-Backend/database"
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

func NewValidator() *validator.Validate {
	return validator.New()
}

func initializeDependency(config *config.Config) (*routes.Dependency, error) {

	wire.Build(
		database.NewDB,
		middlewares.NewMiddlewares,
		NewValidator
		routes.NewDependency,
		userSet,
	)

	return nil, nil
}
