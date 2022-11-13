//go:generate oapi-codegen --config types.yaml ../../ipfs-pinning-service.yaml
//go:generate oapi-codegen --config server.yaml ../../ipfs-pinning-service.yaml

package api

import (
	"context"
	shell "github.com/ipfs/go-ipfs-api"
)

type PinningServer struct {
}

func NewPinningServer() *PinningServer {
	return &PinningServer{}
}

func Reason(reason string) BadRequestJSONResponse {
	return BadRequestJSONResponse{Error: struct {
		Details *string `json:"details,omitempty"`
		Reason  string  `json:"reason"`
	}(struct {
		Details *string
		Reason  string
	}{Details: nil, Reason: reason})}
}

func (p PinningServer) GetPins(ctx context.Context, request GetPinsRequestObject) (GetPinsResponseObject, error) {
	var results PinResults

	sh := shell.NewShell("localhost:5001")

	if request.Params.Cid == nil {
		return GetPins400JSONResponse{Reason("cid is required")}, nil
	}

	for _, cid := range *request.Params.Cid {
		err := sh.Pin(cid)
		if err != nil {
			panic(err)
		}
	}
	return GetPins200JSONResponse(results), nil
}

func (p PinningServer) AddPin(ctx context.Context, request AddPinRequestObject) (AddPinResponseObject, error) {
	//TODO implement me
	return AddPin202JSONResponse{}, nil
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
