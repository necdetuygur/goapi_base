package main

import socketio "github.com/googollee/go-socket.io"

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
		for _, socket := range Sockets {
			socket.Emit("Dream", msg)
		}
	})
	return server
}
