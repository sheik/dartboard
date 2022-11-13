package main

import (
	"dartboard/internal/api"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("dartboard starting")

	swagger, err := api.GetSwagger()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot get swagger")
	}
	swagger.Servers = nil

	pinningServer := api.NewPinningServer()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:Authorization",
		Validator: func(key string, c echo.Context) (bool, error) {
			log.Info().Str("Authorization", key).Msg("request")
			return key == "secret", nil
		},
	}))

	strictHandler := api.NewStrictHandler(pinningServer, nil)
	api.RegisterHandlers(e, strictHandler)

	log.Fatal().Err(e.Start(":8888")).Msg("server exited")
}
