package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

const TestingAddress = "0.0.0.0:8813"

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	//	Format: "progress", // can define default values
}

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
	zerolog.SetGlobalLevel(zerolog.Disabled)
}

func TestMain(m *testing.M) {
	pflag.Parse()
	opts.Paths = pflag.Args()

	status := godog.TestSuite{
		Name: "godogs",
		//		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()

	os.Exit(status)
}

func aListOfPinsIsRequested(ctx context.Context) (context.Context, error) {
	url := fmt.Sprintf("http://%s/pins", ctx.Value("address").(string))
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Authorization", "Bearer secret")
	resp, err := client.Do(req)
	if err != nil {
		return ctx, err
	}
	newCtx := context.WithValue(ctx, "response", resp)
	return newCtx, nil
}

func aPinningServer(ctx context.Context) (context.Context, error) {
	server := echo.New()
	server.HideBanner = true
	server.HidePort = true

	// TODO implement actual auth system
	server.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:Authorization",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == "secret", nil
		},
	}))

	strictHandler := NewStrictHandler(NewPinningServer(), nil)
	RegisterHandlers(server, strictHandler)

	errChan := make(chan error, 1)
	go func() {
		errChan <- server.Start(TestingAddress)
	}()
	select {
	case err := <-errChan:
		return ctx, err
	case <-time.After(200 * time.Millisecond):
		newCtx := context.WithValue(context.WithValue(ctx, "server", server), "address", TestingAddress)
		return newCtx, nil
	}
}

func thereShouldBePinsInResult(ctx context.Context, expectedCount int) (context.Context, error) {
	response := ctx.Value("response").(*http.Response)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return ctx, err
	}
	pin := PinResults{}
	err = json.Unmarshal(body, &pin)
	if err != nil {
		return ctx, err
	}
	if int(pin.Count) != expectedCount {
		return ctx, fmt.Errorf("invalid pin count in response: expected %d, received %d", expectedCount, pin.Count)
	}
	return ctx, nil
}

func aPinIsCreatedWithNameAndCID(ctx context.Context, name, CID string) (context.Context, error) {
	posturl := fmt.Sprintf("http://%s/pins", ctx.Value("address").(string))
	client := http.Client{}
	pin := Pin{
		Cid:  CID,
		Name: &name,
	}
	body, err := json.Marshal(pin)
	if err != nil {
		return ctx, err
	}
	req, err := http.NewRequest(http.MethodPost, posturl, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer secret")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return ctx, err
	}
	newCtx := context.WithValue(ctx, "response", resp)

	return newCtx, nil
}

func anEmptyDatabase(ctx context.Context) (context.Context, error) {
	db, err := sql.Open("sqlite", "pins.sqlite")
	if err != nil {
		return ctx, err
	}
	defer db.Close()
	_, err = db.Exec("drop table if exists pins")
	if err != nil {
		return ctx, err
	}
	_, err = db.Exec(CreateTableSQL)
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}

func theResponseCodeShouldBe(ctx context.Context, expectedCode int) error {
	response := ctx.Value("response").(*http.Response)
	if response.StatusCode != expectedCode {
		return fmt.Errorf("invalid response code: expected %d, recevied %d", expectedCode, response.StatusCode)
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		server := ctx.Value("server").(*echo.Echo)
		if server != nil {
			server.Shutdown(ctx)
		}
		return ctx, nil
	})
	ctx.Step(`^a list of pins is requested$`, aListOfPinsIsRequested)
	ctx.Step(`^a pinning server$`, aPinningServer)
	ctx.Step(`^there should be (\d+) pins in result$`, thereShouldBePinsInResult)
	ctx.Step(`^a pin is created with name "([^"]*)" and CID "([^"]*)"$`, aPinIsCreatedWithNameAndCID)
	ctx.Step(`^an empty database$`, anEmptyDatabase)
	ctx.Step(`^the response code should be (\d+)$`, theResponseCodeShouldBe)
}
