package main

import (
	"fmt"
	"goapi_base/config"
	"goapi_base/model"
	"goapi_base/router"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	socketio "github.com/googollee/go-socket.io"
)

type SocketQueue []socketio.Conn

var Sockets SocketQueue

func SocketServer() *socketio.Server {
	var server = socketio.NewServer(nil)
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		Sockets = append(Sockets, s)
		return nil
	})
	server.OnError("/", func(s socketio.Conn, e error) {})
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {})
	server.OnEvent("/", "Dream", func(s socketio.Conn, msg string) {
		fmt.Println(msg)
		for _, socket := range Sockets {
			socket.Emit("Dream", msg)
		}
	})
	return server
}

func main() {
	// VERİ TABANINDA TABLO YOKSA OLUŞTUR METODLARI
	model.TodoCreateTable()

	// WEB FRAMEWORK TANIMLAMALARI
	ec := echo.New()
	ec.HideBanner = true
	ec.Use(middleware.Logger())
	ec.Use(middleware.Recover())
	ec.Use(middleware.CORS())

	// WEB FRAMEWORK AUTH
	ec.POST("/Login", Login)

	///
	var server = SocketServer()
	go server.Serve()
	defer server.Close()
	ec.Any("/socket.io/", func(c echo.Context) error {
		server.ServeHTTP(c.Response(), c.Request())
		return nil
	})
	///

	e := ec.Group("")
	e.Use(middleware.JWT([]byte(config.AUTH_JWT_TOKEN)))

	// WEB FRAMEWORK YÖNLENDİRİCİLERİ
	router.TodoRouter(e)

	// WEB FRAMEWORK BAŞLAT
	host := "localhost"
	port := config.PORT
	if os.Getenv("PORT") != "" {
		host = "0.0.0.0"
		port = os.Getenv("PORT")
	}
	ec.Logger.Fatal(ec.Start(host + ":" + port))
}
