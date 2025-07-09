package wireinject

import (
	"github.com/devanadindraa/Evermos-Backend/database"
	"github.com/devanadindraa/Evermos-Backend/middlewares"
	"github.com/devanadindraa/Evermos-Backend/routes"
	"github.com/devanadindraa/Evermos-Backend/utils/config"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
)

func initializeDependency(config *config.Config) (*routes.Dependency, error) {

	wire.Build(
		database.NewDB,
		validator.New,
		middlewares.NewMiddlewares,
		routes.NewDependency,
	)

	return nil, nil
}
