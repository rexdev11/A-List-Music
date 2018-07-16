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
	"time"
)

const ServerHost = string("localhost:14591")

type MainOptions struct {
	setUUID string
	consoleCh chan string
}

var consoleCh = make(chan string)

func main() {

	go chConsole()
	go Nexus()

	startupUUIDNum := rand.Float64()
	timestamp := time.Now()
	startupUUID := strconv.FormatFloat(startupUUIDNum, 'f', 25, 64) + timestamp.String()
	fmt.Println(startupUUID)

	options := MainOptions{
		setUUID: startupUUID,
		consoleCh: consoleCh,
	}

	// Nothing should execute past this. //
	StartServer(options)
}

func goC(f *interface{}){
	chConsole()
}

func chConsole() {
	for message := range consoleCh {
		fmt.Printf("%s", message)
	}
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
		HostingData:    getHostData(),
		ServerRoomName: options.setUUID,
		ConsoleCh:      consoleCh,
	}
	serv := server.BuildServer(serverOptions)
	serv.Logger().Debug("on connection")

	err := serv.Run(iris.TLS(
		ServerHost,
		"./alist.cert",
		"./alist.key",
		),
		iris.WithConfiguration(iris.Configuration{ // default configuration:
			DisableStartupLog:                 false,
			DisableInterruptHandler:           false,
			DisablePathCorrection:             false,
			EnablePathEscape:                  false,
			FireMethodNotAllowed:              false,
			DisableBodyConsumptionOnUnmarshal: false,
			DisableAutoFireStatusCode:         false,
			EnableOptimizations: 				true,
			TimeFormat:                        "Mon, 02 Jan 2006 15:04:05 GMT",
			Charset:                           "UTF-8",
	}))

	if err != nil {
		//println(err)
		utilities.ErrorHandler(err)
	}
	fmt.Println("starting server on ", ServerHost)
	target, _ := url.Parse("http://localhost:80")
	host.NewProxy(ServerHost, target).ListenAndServe()
}

// This will coordinate the different modules into routines.
func Nexus() {

	for job := range store.Client().Jobs {

		fmt.Printf("%s", job)
	}

	for file := range server.Client().FileUploaded {
		fmt.Println(file)
	}

	for transcodes := range transcoder.Client().Jobs {
		fmt.Println(transcodes)
	}
}