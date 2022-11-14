//go:generate oapi-codegen --config types.yaml ../../ipfs-pinning-service.yaml
//go:generate oapi-codegen --config server.yaml ../../ipfs-pinning-service.yaml

package api

import (
	"context"
	"database/sql"
	_ "embed"
	"github.com/google/uuid"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/rs/zerolog/log"
	_ "modernc.org/sqlite"
	"strings"
	"time"
)

var (
	//go:embed sql/pin.sql
	CreateTableSQL string
)

type PinningServer struct {
	pins map[string]PinStatus
}

func NewPinningServer() *PinningServer {
	db, err := sql.Open("sqlite", "pins.sqlite")
	if err != nil {
		log.Error().Err(err).Msg("unable to open database")
	}
	defer db.Close()
	_, err = db.Exec(CreateTableSQL)
	if err != nil {
		log.Error().Err(err).Msg("unable to open database")
	}
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
		return GetPins5XXJSONResponse{StatusCode: 500}, nil
	}
	defer db.Close()
	status := "pinned"
	if request.Params.Status != nil {
		var statusList []string
		for _, s := range *request.Params.Status {
			statusList = append(statusList, string(s))
		}
		status = strings.Join(statusList, ", ")
	}
	after := time.Now().UTC().AddDate(-100, 0, 0).Round(time.Second)
	if request.Params.After != nil {
		after = *request.Params.After
	}
	before := time.Now().UTC().AddDate(1, 0, 0).Round(time.Second)
	if request.Params.Before != nil {
		before = *request.Params.Before
	}
	limit := 10
	if request.Params.Limit != nil {
		limit = int(*request.Params.Limit)
	}

	var stmt *sql.Stmt
	var params []interface{}

	if request.Params.Cid != nil {
		stmt, err = db.Prepare("select cid, name, request_id, created_at, status from pins where cid IN (?)")
		for _, cid := range *request.Params.Cid {
			params = append(params, cid)
		}
	} else {
		stmt, err = db.Prepare("select cid, name, request_id, created_at, status from pins where status IN (?) and created_at < ? and created_at > ? limit ?")
		params = append(params, status)
		params = append(params, before)
		params = append(params, after)
		params = append(params, limit)
	}
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return GetPins5XXJSONResponse{StatusCode: 500}, nil
	}
	rows, err := stmt.Query(params...)
	if err != nil {
		log.Error().Err(err).Msg("unable to select from pins")
		return GetPins5XXJSONResponse{StatusCode: 500}, nil
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
		return GetPins5XXJSONResponse{StatusCode: 500}, nil
	}

	results.Count = int32(len(results.Results))

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
		return AddPin400JSONResponse{Reason("unable to pin")}, nil
	}

	requestId, err := uuid.NewUUID()
	if err != nil {
		return AddPin5XXJSONResponse{StatusCode: 500}, nil
	}

	result := PinStatus{
		Created:   time.Now().UTC().Round(time.Second),
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

	db, err := sql.Open("sqlite", "pins.sqlite")
	if err != nil {
		log.Error().Err(err).Msg("unable to open database")
		return AddPin5XXJSONResponse{StatusCode: 500}, nil
	}
	defer db.Close()

	stmt, err := db.Prepare("insert into pins (cid, name, request_id, created_at, status) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare database query")
		return AddPin5XXJSONResponse{StatusCode: 500}, nil
	}

	_, err = stmt.Exec(result.Pin.Cid, result.Pin.Name, result.Requestid, result.Created, "pinned")
	if err != nil {
		log.Error().Err(err).Msg("unable to execute database query")
		return AddPin5XXJSONResponse{StatusCode: 500}, nil
	}

	return AddPin202JSONResponse(result), nil
}

func (p PinningServer) DeletePinByRequestId(ctx context.Context, request DeletePinByRequestIdRequestObject) (DeletePinByRequestIdResponseObject, error) {
	if request.Requestid == "" {
		return DeletePinByRequestId400JSONResponse{Reason("requestId is required")}, nil
	}

	sh := shell.NewShell("localhost:5001")

	db, err := sql.Open("sqlite", "pins.sqlite")
	if err != nil {
		log.Error().Err(err).Msg("unable to open database")
		return DeletePinByRequestId5XXJSONResponse{StatusCode: 500}, nil
	}
	defer db.Close()

	stmt, err := db.Prepare("select request_id, cid from pins where request_id = ?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return DeletePinByRequestId5XXJSONResponse{StatusCode: 500}, nil
	}

	rows, err := stmt.Query(request.Requestid)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query")
		return DeletePinByRequestId5XXJSONResponse{StatusCode: 500}, nil
	}
	defer rows.Close()

	for rows.Next() {
		var requestId string
		var cid string

		rows.Scan(&requestId, &cid)

		// TODO pins should only be unpinned if this is not pinned over all accounts
		err := sh.Unpin(cid)
		if err != nil {
			log.Error().Str("CID", cid).Str("RequestId", requestId).Err(err).Msg("could not unpin cid")
		}
	}

	err = rows.Err()
	if err != nil {
		log.Error().Err(err).Msg("error fetching rows")
		return DeletePinByRequestId5XXJSONResponse{StatusCode: 500}, nil
	}

	stmt, err = db.Prepare("delete from pins where request_id = ?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare delete query")
		return DeletePinByRequestId5XXJSONResponse{StatusCode: 500}, nil
	}
	_, err = stmt.Exec(request.Requestid)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute delete query")
		return DeletePinByRequestId5XXJSONResponse{StatusCode: 500}, nil
	}
	return DeletePinByRequestId202Response{}, nil
}

func (p PinningServer) GetPinByRequestId(ctx context.Context, request GetPinByRequestIdRequestObject) (GetPinByRequestIdResponseObject, error) {
	db, err := sql.Open("sqlite", "pins.sqlite")
	if err != nil {
		log.Error().Err(err).Msg("unable to open database")
		return GetPinByRequestId5XXJSONResponse{StatusCode: 500}, nil
	}
	defer db.Close()

	stmt, err := db.Prepare("select cid, name, request_id, created_at, status from pins where request_id = ?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return GetPinByRequestId5XXJSONResponse{StatusCode: 500}, nil
	}

	rows, err := stmt.Query(request.Requestid)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query")
		return GetPinByRequestId5XXJSONResponse{StatusCode: 500}, nil
	}
	defer rows.Close()

	pinStatus := PinStatus{}
	for rows.Next() {
		rows.Scan(&pinStatus.Pin.Cid, &pinStatus.Pin.Name, &pinStatus.Requestid, &pinStatus.Created, &pinStatus.PinningStatus)
	}
	err = rows.Err()
	if err != nil {
		log.Error().Err(err).Msg("error fetching rows")
		return GetPinByRequestId5XXJSONResponse{StatusCode: 500}, nil
	}

	return GetPinByRequestId200JSONResponse(pinStatus), nil
}

func (p PinningServer) ReplacePinByRequestId(ctx context.Context, request ReplacePinByRequestIdRequestObject) (ReplacePinByRequestIdResponseObject, error) {
	//TODO implement ReplacePinByRequestId
	panic("implement me")
}
