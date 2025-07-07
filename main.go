package main

import (
	"github.com/devanadindraa/Evermos-Backend/utils/config"
	"github.com/go-playground/validator/v10"
)

func main() {
	conf := config.NewConfig()

	err := validator.New().Struct(conf)
	if err != nil {
		panic(err)
	}

}
