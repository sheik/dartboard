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
		Command: "go build -o . ./...",
		Depends: "clean api test",
		Default: true,
		Help:    "build binaries",
	},
	"test": {
		Command: "go test ./...",
		Depends: "clean",
		Help:    "run bdd tests",
	},
}

func main() {
	fab.Run(plan)
}
