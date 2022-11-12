//go:generate oapi-codegen --config types.yaml ../../ipfs-pinning-service.yaml
//go:generate oapi-codegen --config server.yaml ../../ipfs-pinning-service.yaml

package api

import (
	"context"
)

type PinningServer struct {
}

func NewPinningServer() *PinningServer {
	return &PinningServer{}
}

func (p PinningServer) GetPins(ctx context.Context, request GetPinsRequestObject) (GetPinsResponseObject, error) {
	var results PinResults
	return GetPins200JSONResponse(results), nil
}

func (p PinningServer) AddPin(ctx context.Context, request AddPinRequestObject) (AddPinResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (p PinningServer) DeletePinByRequestId(ctx context.Context, request DeletePinByRequestIdRequestObject) (DeletePinByRequestIdResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (p PinningServer) GetPinByRequestId(ctx context.Context, request GetPinByRequestIdRequestObject) (GetPinByRequestIdResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (p PinningServer) ReplacePinByRequestId(ctx context.Context, request ReplacePinByRequestIdRequestObject) (ReplacePinByRequestIdResponseObject, error) {
	//TODO implement me
	panic("implement me")
}
