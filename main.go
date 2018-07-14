package main

import (
	"fmt"
	"a-list-music/server"
	"github.com/kataras/iris"
	"a-list-music/store"
	"a-list-music/transcoder"
	"a-list-music/utilities"
	"net/url"
	"github.com/kataras/iris/core/host"
)

const ServerHost = string("localhost:10911")

func main() {

	StartServer()
	Nexus()
}

func StartServer() {
	fmt.Println("starting Server")
		serv := server.BuildServer(server.Data{ ServerHostPort: ServerHost })
		serv.Logger().Info("on connection")
		fmt.Println("starting server on ", ServerHost)
		target, _ := url.Parse("localhost:*")
		go host.NewProxy(ServerHost, target).ListenAndServe()

	err := serv.Run(iris.TLS(ServerHost, "./alist.cert", "./alist.key"), iris.WithConfiguration(iris.Configuration{ // default configuration:
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