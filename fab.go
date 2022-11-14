//go:build exclude
// +build exclude

package main

import "github.com/sheik/fab"

const imageName = "dartboard-util:latest"

var (
	image   = fab.Container(imageName).Mount("$PWD", "/app")
	version = fab.GetVersion()
	uid     = fab.Output("id -u")
	gid     = fab.Output("id -g")
)

var plan = fab.Plan{
	"dartboard-util": {
		Command: "docker build -f docker/dockerfiles/dartboard-util/Dockerfile . -t dartboard-util:latest",
		Check:   fab.ImageExists(imageName),
	},
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
	"deb": {
		Command: image.Run("fpm -s dir -t deb -n dartboard -v %s usr", version),
		Depends: "dartboard-util",
	},
	"rpm": {
		Command: image.Run("fpm -s dir -t rpm -n dartboard -v %s usr", version),
		Depends: "dartboard-util",
	},
	"package": {
		Command: image.Run("chown %s:%s *.deb *.rpm", uid, gid),
		Depends: "deb rpm",
	},
}

func main() {
	fab.Run(plan)
}
