package main

import (
	"goapi_base/config"
	"goapi_base/model"
	"goapi_base/router"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

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

	// SOCKET.IO
	var socketServer = SocketServer()
	go socketServer.Serve()
	defer socketServer.Close()
	ec.Any("/socket.io/", func(c echo.Context) error {
		socketServer.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	// AUTH_JWT_TOKEN
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
