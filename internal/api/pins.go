//go:generate oapi-codegen --config types.yaml ../../ipfs-pinning-service.yaml
//go:generate oapi-codegen --config server.yaml ../../ipfs-pinning-service.yaml

package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type PinningServer struct {
}

var _ ServerInterface = (*PinningServer)(nil)

func NewPinningServer() *PinningServer {
	return &PinningServer{}
}

// List pin objects
// (GET /pins)
func (ps *PinningServer) GetPins(ctx echo.Context, params GetPinsParams) error {
	return ctx.JSON(http.StatusOK, nil)
}

// Add pin object
// (POST /pins)
func (ps *PinningServer) AddPin(ctx echo.Context) error {
	return nil
}

// Remove pin object
// (DELETE /pins/{requestid})
func (ps *PinningServer) DeletePinByRequestId(ctx echo.Context, requestid string) error {
	return nil
}

// Get pin object
// (GET /pins/{requestid})
func (ps *PinningServer) GetPinByRequestId(ctx echo.Context, requestid string) error {
	return nil
}

// Replace pin object
// (POST /pins/{requestid})
func (ps *PinningServer) ReplacePinByRequestId(ctx echo.Context, requestid string) error {
	return nil
}
