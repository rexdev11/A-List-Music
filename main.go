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
	"strings"
	"github.com/tensoflow/tensorflow/tensorflow/go"
)

// Main Server Settings
var Host_Data = func() server.HostingInfo {
	const Host = string("localhost")
	const Port = 8080
	const Protocol = "https"

	return server.HostingInfo{
		Paths: server.HostPaths{
			SocketMainRoomName: ":Main",
			Path:               string(Host + ":" + strconv.Itoa(Port)),
			Host:               Host,
			Protocol:           Protocol,
			URI:                Protocol + "://" + Host + ":" + strconv.Itoa(Port),
			Port:               Port,
		},
	}
}()

type ServerOptions struct {
	serverRndUUID 	string
	consoleCh 		chan string
}

var consoleCh = make(chan string)

var rdmUUID = func() string {
	startupUUIDNum := rand.Float64()
	timestamp := time.Now()
	return strings.Trim(strconv.FormatFloat(startupUUIDNum, 'f', 25, 64) + timestamp.String(), " ")
}()

func main() {
	go chConsole()
	go Nexus()

	fmt.Println(Host_Data)

	options := ServerOptions{
		serverRndUUID: rdmUUID,
		consoleCh:     consoleCh,
	}

	// Nothing will execute past StartServer.
	// (as long as the server loop doesn't bork...)
	StartServer(options)
}

func chConsole() {
	for message := range consoleCh {
		fmt.Printf("%s", message)
	}
}

func StartServer(options ServerOptions) {
	fmt.Println("starting Server")
	serverOptions := server.ServerOptions {
		HostingData:    Host_Data,
		ServerRoomName: options.serverRndUUID,
		ConsoleCh:      consoleCh,
	}

	serv := server.BuildServer(serverOptions)
	serv.Logger().Debug("on connection")

	err := serv.Run(iris.TLS(
		Host_Data.Paths.Path,
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
			EnableOptimizations:				true,
			TimeFormat:                        "Mon, 02 Jan 2006 15:04:05 GMT",
			Charset:                           "UTF-8",
	}))

	if err != nil {

		//println(err)
		utilities.ErrorHandler(err)

	}
	fmt.Println("starting server on ", Host_Data.Paths.Host)
	target, _ := url.Parse("localhost:80")
	host.NewProxy(Host_Data.Paths.URI, target).ListenAndServe()
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