package server

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/websocket"
	"sync/atomic"
	"fmt"
	"github.com/kataras/iris/context"
	"a-list-music/utilities"
	"net/http"
	"html/template"
	"encoding/json"
)

var fileUploadedChan chan utilities.Action
var Client = func() AListServerClient {
	return AListServerClient{
		FileUploaded: fileUploadedChan,
	}
}

type AListServer interface {
	BuildServer() (server iris.Application)
}

type AListServerClient struct {
	FileUploaded chan utilities.Action
}

type HostPaths struct {
	Name     string
	Path     string
	Host     string
	Protocol string
	URI      string
	Port     int
}
type HostingInfo struct {
	Paths HostPaths
}

type ServerOptions struct {
	HostingData HostingInfo
	StartUpUUID string
}

func BuildServer(options ServerOptions) (server *iris.Application) {
	app := iris.Default()
	app.RegisterView(iris.HTML("./views", ".html"))
	app.Get("/", func(ctx iris.Context) {
		var body string

		ctx.ResponseWriter()
		ctx.ResetResponseWriter(ctx.ResponseWriter())

		// options to JSON
		jsonData, err := json.Marshal(options)
		utilities.ErrorHandler(err)

		// Process Template
		tmplt, err := template.ParseFiles("./views/index.html")

		// set JSON
		err = tmplt.Execute(ctx.ResponseWriter(), template.HTML(jsonData))
		utilities.ErrorHandler(err)

		ctx.WriteString(body)
		ctx.ViewLayout(body)
	})

	app.Get("/style-sheet", func(ctx context.Context) {
		ctx.ServeFile("./views/main.style.css", false)
	})

	app.Get("/alist-service", func(ctx context.Context) {
		ctx.ServeFile("./workers/alist-service.worker.js", false)
	})

	app.Get("/admin", func(ctx context.Context) {
		ctx.View("admin.html")
	})

	mvc.Configure(app.Party("/websocket"), configureMVC)

	return app
}

func configureMVC(m *mvc.Application) {
	ws := websocket.New(websocket.Config{
		CheckOrigin: func(r *http.Request) bool {
			fmt.Println(r)
			return true
		},
		IDGenerator: func(ctx context.Context) string {
			var count= int(0)
			var name= "ClientID" + string(count+1)
			fmt.Println(name)
			return name
		},
	})
	// http://localhost:8080/websocket/iris-ws.js
	m.Router.Any("/iris-ws.js", websocket.ClientHandler())

	// This will bind the result of ws.Upgrade which is a websocket.Connection
	// to the controller(s) served by the `m.Handle`.
	m.Register(ws.Upgrade)
	m.Handle(new(websocketController))
}

var visits uint64
var roomPrefix string

func increment() uint64 {
	return atomic.AddUint64(&visits, 1)
}

func decrement() uint64 {
	return atomic.AddUint64(&visits, ^uint64(0))
}

type websocketController struct {
	Conn websocket.Connection
}

func (c *websocketController) onLeave(roomName string) {
	newCount := decrement()
	c.Conn.To(websocket.Broadcast).Emit("visit", newCount)
}

func (c *websocketController) update() {
	newCount := increment()
	c.Conn.To(websocket.All).Emit("visit", newCount)
}

func (c *websocketController) Get( startUUID string /* websocket.Connection could be lived here as well, it doesn't matter */ ) {
	roomPrefix = startUUID
	c.Conn.Join(roomPrefix + ":Main")
	c.Conn.OnLeave(c.onLeave)
	c.Conn.On("visit", c.update)
	c.fileSockets(roomPrefix)
	c.Conn.Wait()
}

func(c *websocketController) fileSockets(roomPrefix string) {
	var RoomName = roomPrefix + "file_upload"
	c.Conn.Join(RoomName)
	if c.Conn.IsJoined(RoomName) {
		var onUpload = func(data []byte) {
			fmt.Println("file received", data)
		}
		c.Conn.On("upload", onUpload)
	}

	c.Conn.Emit("FileUpload::Done", nil)
}