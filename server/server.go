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
	"path"
)

var MainServerRoom string
var fileUploadedChan chan utilities.Action
var consoleCh chan string
var ClieSnt = func() AListServerClient {
	return AListServerClient {
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
	SocketMainRoomName string
	Path               string
	Host               string
	Protocol           string
	URI                string
	Port               int
}
type HostingInfo struct {
	Paths HostPaths
}

type ServerOptions struct {
	HostingData HostingInfo
	ServerRoomName string
	ConsoleCh      chan string
}

type SiteData struct {
	HostingData HostingInfo
	ServerRoomName string
}

var body string

func BuildServer(options ServerOptions) (server *iris.Application) {
	consoleCh = options.ConsoleCh
	app := iris.Default()

	MainServerRoom = options.ServerRoomName + ":Main"

	app.RegisterView(iris.HTML("./views", ".html"))

	app.Get("/", func(ctx iris.Context) {
		consoleCh <- "foo"
		res := ctx.ResponseWriter()
		origin := res.Header().Get("origin")
		consoleCh <- string("I ma the 70 origin" + origin)
		res.Push(origin, nil)
		data := SiteData{options.HostingData, options.ServerRoomName}
		ctx.ResponseWriter()
		ctx.ResetResponseWriter(ctx.ResponseWriter())

		// options to JSON
		jsonData, err := json.Marshal(data)

		utilities.ErrorHandler(err)

		// Process Template
		tmplt, err := template.ParseFiles("./views/index.html")

		// set JSON
		err = tmplt.Execute(ctx.ResponseWriter(), template.HTML(jsonData))
		utilities.ErrorHandler(err)

		// Respond
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
		res := ctx.ResponseWriter()
		origin := res.Header().Get("origin")
		consoleCh <- origin

		res.Push(origin, nil )

		data := SiteData{options.HostingData, options.ServerRoomName}
		ctx.ResponseWriter()
		ctx.ResetResponseWriter(ctx.ResponseWriter())

		// options to JSON
		jsonData, err := json.Marshal(data)

		utilities.ErrorHandler(err)

		// Process Template
		tmplt, err := template.ParseFiles("./views/admin.html")

		// set JSON
		err = tmplt.Execute(ctx.ResponseWriter(), template.HTML(jsonData))
		utilities.ErrorHandler(err)

		ctx.WriteString(body)
		ctx.ViewLayout(body)
	})

	options.ConsoleCh <- "Sockets Initializing"

	app.StaticWeb("/js/", path.Join(utilities.CWD(), "/web-src/js/") )

	mvc.Configure(app.Party("/websocket"), configureMVC)

	return app
}

func configureMVC(m *mvc.Application) {
	ws := websocket.New(websocket.Config{
		CheckOrigin: func(r *http.Request) bool {
			consoleCh <- "Checking ORIGIN"
			return true
		},
		IDGenerator: func(ctx context.Context) string {
			var count int
			var name = "ClientID" + string(count+1)
			consoleCh <- "Client SocketMainRoomName"
			consoleCh <- name
			consoleCh <- ctx.String()
			return name
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	})
	// http://localhost:8080/websocket/iris-ws.js
	m.Router.Any("/iris-ws.js", websocket.ClientHandler())

	// This will bind the result of ws.Upgrade which is a websocket.Connection
	// to the controller(s) served by the `m.Handle`.
	m.Register(ws.Upgrade)
	m.Handle(new(websocketController))
	//ws.Join(MainServerRoom)
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
	c.Conn.OnLeave(c.onLeave)
	c.Conn.On("visit", c.update)
	c.Conn.Wait()

	fmt.Println("Sockets Waiting")
	var RoomName = roomPrefix + "file_upload"
	c.Conn.Join(RoomName)
	if c.Conn.IsJoined(RoomName) {
		var onUpload = func(data []byte) {
			fmt.Println("file received", data)
		}
		c.Conn.On("upload", onUpload)
	}

	c.Conn.Emit("FileUpload::Done", nil)
	c.Conn.Wait()
}