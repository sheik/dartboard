package server

import (
	"dartboard/internal/api"
	"fmt"
	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"net/http"
)

func CreateMiddleware(v JWSValidator) (func(http.Handler) http.Handler, error) {
	spec, err := api.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("loading spec: %w", err)
	}

	validator := middleware.OapiRequestValidatorWithOptions(spec,
		&middleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: NewAuthenticator(v),
			},
		})

	return validator, nil
}
