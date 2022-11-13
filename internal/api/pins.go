//go:generate oapi-codegen --config types.yaml ../../ipfs-pinning-service.yaml
//go:generate oapi-codegen --config server.yaml ../../ipfs-pinning-service.yaml

package api

import (
	"context"
	"github.com/google/uuid"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/rs/zerolog/log"
	"time"
)

type PinningServer struct {
	pins map[string]PinStatus
}

func NewPinningServer() *PinningServer {
	return &PinningServer{pins: make(map[string]PinStatus)}
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

	pins, err := sh.Pins()
	if err != nil {
		return GetPins5XXJSONResponse{}, nil
	}

	for cid, _ := range pins {
		results.Results = append(results.Results, PinStatus{
			Created: time.Time{},
			Pin: Pin{
				Cid: cid,
			},
			Requestid:     "",
			PinningStatus: "pinned",
		})
	}

	return GetPins200JSONResponse(results), nil
}

func (p PinningServer) AddPin(ctx context.Context, request AddPinRequestObject) (AddPinResponseObject, error) {
	sh := shell.NewShell("localhost:5001")
	if request.Body.Cid == "" {
		return AddPin400JSONResponse{Reason("cid is required")}, nil
	}

	err := sh.Pin(request.Body.Cid)
	if err != nil {
		log.Error().Err(err).Msg("unable to pin")
		return AddPin400JSONResponse{Reason(err.Error())}, nil
	}

	requestId, err := uuid.NewUUID()
	if err != nil {
		return AddPin5XXJSONResponse{}, nil
	}

	result := PinStatus{
		Created:   time.Now().UTC(),
		Delegates: nil,
		Info:      nil,
		Pin: Pin{
			Cid:     request.Body.Cid,
			Meta:    request.Body.Meta,
			Name:    request.Body.Name,
			Origins: request.Body.Origins,
		},
		Requestid:     requestId.String(),
		PinningStatus: "pinned",
	}

	p.pins[requestId.String()] = result

	return AddPin202JSONResponse(result), nil
}

func (p PinningServer) DeletePinByRequestId(ctx context.Context, request DeletePinByRequestIdRequestObject) (DeletePinByRequestIdResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (p PinningServer) GetPinByRequestId(ctx context.Context, request GetPinByRequestIdRequestObject) (GetPinByRequestIdResponseObject, error) {
	pinStatus, ok := p.pins[request.Requestid]
	if !ok {
		return GetPinByRequestId404JSONResponse{}, nil
	}

	return GetPinByRequestId200JSONResponse(pinStatus), nil
}

func (p PinningServer) ReplacePinByRequestId(ctx context.Context, request ReplacePinByRequestIdRequestObject) (ReplacePinByRequestIdResponseObject, error) {
	//TODO implement me
	panic("implement me")
}
