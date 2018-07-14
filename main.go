package main

import (
	"fmt"
	"a-list-music/server"
	"github.com/kataras/iris"
	"a-list-music/store"
	"a-list-music/transcoder"
)

const ServerHost = string("localhost:9033")

func main() {
	var st = string("calling AListTranscoder")
	fmt.Println(st)
	StartServer()
}

func StartServer() {
	fmt.Println("starting Server")
	serv := server.BuildServer()
	serv.Run(iris.Addr(ServerHost))
	Nexus()
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