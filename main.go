package main

import (
			"fmt"
	"github.com/kataras/iris"
	"a-list/server"
	"a-list/transcoder"
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

	go transcoder.init()
	// main thread ends with Server.Run
	StartServer()

}

func StartServer() {
	fmt.Println("starting Server")
	server := server.BuildServer()
	server.Run(iris.Addr("localhost:2820"))
}

