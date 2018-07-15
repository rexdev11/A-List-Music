package main

import (
	"fmt"
	"a-list-music/server"
	"a-list-music/store"
	"a-list-music/transcoder"
	"net/url"
	"github.com/kataras/iris/core/host"
	"github.com/kataras/iris"
	"a-list-music/utilities"
	"math/rand"
	"strconv"
)

const ServerHost = string("localhost:14121")

type MainOptions struct {
	setUUID string
}

func main() {
	startupUUIDNum := rand.Float64()
	startupUUID := strconv.FormatFloat(startupUUIDNum, 'f', 6, 64)

	options := MainOptions{
		setUUID: startupUUID,
	}

	StartServer(options)
	Nexus()
}

func getHostData() server.HostingInfo {
	hostPaths := server.HostingInfo{
		Paths: server.HostPaths{
			Name:     "MainServer",
			Path:     "localhost:12121",
			Host:     "localhost",
			Protocol: "https",
			URI:      "https://localhost:12121",
			Port:     134545,
		},
	}
	return hostPaths
}

func StartServer(options MainOptions) {
	fmt.Println("starting Server")
	serverOptions := server.ServerOptions{
		HostingData: getHostData(),
		StartUpUUID: options.setUUID,
	}
	serv := server.BuildServer(serverOptions)
	serv.Logger().Info("on connection")
	fmt.Println("starting server on ", ServerHost)
	target, _ := url.Parse("localhost:433")

	go host.NewProxy(ServerHost, target).ListenAndServe()
		err := serv.Run(iris.TLS(
			ServerHost,
			"./alist.cert",
			"./alist.key",
			), iris.WithConfiguration(
				iris.Configuration{ // default configuration:
				DisableStartupLog:                 false,
				DisableInterruptHandler:           false,
				DisablePathCorrection:             false,
				EnablePathEscape:                  false,
				FireMethodNotAllowed:              false,
				DisableBodyConsumptionOnUnmarshal: false,
				DisableAutoFireStatusCode:         false,
				TimeFormat:                        "Mon, 02 Jan 2006 15:04:05 GMT",
				Charset:                           "UTF-8",
			}))
	if err != nil {
		utilities.ErrorHandler(err)
	}
}

// This will coordinate the different modules into routines.
func Nexus() {
	for job := range store.Client().Jobs {
		fmt.Println(job)
	}

	for file := range server.Client().FileUploaded {
		fmt.Println(file)
	}

	for transcodes := range transcoder.Client().Jobs {
		fmt.Println(transcodes)
	}
}