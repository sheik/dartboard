//go:generate oapi-codegen --config types.yaml ../../ipfs-pinning-service.yaml
//go:generate oapi-codegen --config server.yaml ../../ipfs-pinning-service.yaml

package api

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/rs/zerolog/log"
	_ "modernc.org/sqlite"
	"strings"
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

	db, err := sql.Open("sqlite", "pins.sqlite")
	if err != nil {
		log.Error().Err(err).Msg("unable to open database")
		return GetPins5XXJSONResponse{}, nil
	}
	status := "pinned"
	if request.Params.Status != nil {
		var statusList []string
		for _, s := range *request.Params.Status {
			statusList = append(statusList, string(s))
		}
		status = strings.Join(statusList, ", ")
	}
	after := time.Now().AddDate(-100, 0, 0)
	if request.Params.After != nil {
		after = *request.Params.After
	}
	before := time.Now().AddDate(1, 0, 0)
	if request.Params.Before != nil {
		before = *request.Params.Before
	}
	limit := 10
	if request.Params.Limit != nil {
		limit = int(*request.Params.Limit)
	}
	stmt, err := db.Prepare("select cid, name, request_id, created_at, status from pins where status IN (?) and created_at < ? and created_at > ? limit ?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return GetPins5XXJSONResponse{}, nil
	}

	rows, err := stmt.Query(status, before, after, limit)
	if err != nil {
		log.Error().Err(err).Msg("unable to select from pins")
		return GetPins5XXJSONResponse{}, nil
	}
	defer rows.Close()

	for rows.Next() {
		pin := PinStatus{}
		rows.Scan(&pin.Pin.Cid, &pin.Pin.Name, &pin.Requestid, &pin.Created, &pin.PinningStatus)
		results.Results = append(results.Results, pin)
	}
	err = rows.Err()
	if err != nil {
		log.Error().Err(err).Msg("error fetching rows")
		return GetPins5XXJSONResponse{}, nil
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
	db, err := sql.Open("sqlite", "pins.sqlite")
	if err != nil {
		log.Error().Err(err).Msg("unable to open database")
		return AddPin5XXJSONResponse{}, nil
	}

	stmt, err := db.Prepare("insert into pins (cid, name, request_id, created_at, status) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare database query")
		return AddPin5XXJSONResponse{}, nil
	}

	_, err = stmt.Exec(result.Pin.Cid, result.Pin.Name, result.Requestid, result.Created, "pinned")
	if err != nil {
		log.Error().Err(err).Msg("unable to execute database query")
		return AddPin5XXJSONResponse{}, nil
	}

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
