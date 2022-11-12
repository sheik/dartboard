package main

import (
	"dartboard/internal/api"
	"dartboard/internal/server"
	"fmt"
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

	authenticator, err := server.NewFakeAuthenticator()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create authenticator")
	}
	/*
		mw, err := server.CreateMiddleware(authenticator)
		if err != nil {
			log.Fatal().Err(err).Msg("unable to create auth middleware")
		}
	*/

	e.Use(middleware.Logger())
	//	e.Use(echo.WrapMiddleware(mw))

	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:x-api-key",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == "secret", nil
		},
	}))

	api.RegisterHandlers(e, pinningServer)

	fmt.Println(e.Routers())

	readerJWS, err := authenticator.CreateJWSWithClaims([]string{})
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create reader JWS")
	}

	writerJWS, err := authenticator.CreateJWSWithClaims([]string{"things:w"})
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create writer JWS")
	}

	log.Info().Str("token", string(readerJWS)).Msg("reader token")
	log.Info().Str("token", string(writerJWS)).Msg("writer token")

	log.Fatal().Err(e.Start(":8888")).Msg("server exited")
}
