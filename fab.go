//go:build exclude
// +build exclude

package main

import "github.com/sheik/fab"

var plan = fab.Plan{
	"clean": {
		Command: "rm -rf $(ls cmd)",
		Help:    "clean build artifacts from repo",
	},
	"oapi-codegen": {
		Command: "go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest",
	},
	"api": {
		Command: "go generate ./...",
		Help:    "generate api helpers from swagger",
		Depends: "oapi-codegen",
	},
	"build": {
		Command: "CGO_ENABLED=0 go build -o . ./...",
		Depends: "clean api",
		Default: true,
		Help:    "build binaries",
	},
	"test": {
		Command: "go test ./... -v",
		Depends: "clean",
		Help:    "run bdd tests (IPFS must be running locally on port 5001)",
		Gate:    fab.Exec("lsof -i -P -n | grep LISTEN | grep 5001 2>&1 > /dev/null"),
	},
	"docker-image": {
		Command: "docker build . -t dartboard:latest",
		Depends: "build",
		Help:    "build a docker image of dartboard",
	},
}

func main() {
	fab.Run(plan)
}
