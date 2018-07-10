package main

import (
		"a-list/transcoder"
	"fmt"
	"github.com/kataras/iris"
	"a-list/server"
	)
var shutdownServer bool

func main() {
	//soundFile, err := os.Open("./sound-files/dummy.wav")
	//if err != nil {
	//	panic(err)
	//}
	shutdownServer = false
	fmt.Println("calling Transcoder")

	// initialize transcoder channel

	go transcoder.Transcode()

	// main thread ends with Server.Run
	StartServer()

}

func StartServer() {
	fmt.Println("starting Server")
	server := server.BuildServer()
	server.Run(iris.Addr("localhost:2820"))
}

