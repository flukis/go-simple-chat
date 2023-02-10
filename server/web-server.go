package main

import (
	"fmt"
	"io"
	"sync"

	"golang.org/x/net/websocket"
)

type Server struct {
	mu    *sync.Mutex
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) HandleBroadcast(ws *websocket.Conn) {
	fmt.Printf("new broadcasting from %s\n", ws.RemoteAddr())
	s.conns[ws] = true
	func(ws *websocket.Conn) {
		buff := make([]byte, 1024)
		for {
			n, err := ws.Read(buff)
			if err != nil {
				if err == io.EOF {
					fmt.Printf("error end of file, %s\n", err.Error())
					break
				}

				fmt.Printf("read error, %s\n", err.Error())
				continue
			}

			msg := buff[:n]
			for ws := range s.conns {
				go func(ws *websocket.Conn) {
					if _, err := ws.Write(msg); err != nil {
						fmt.Printf("write error, %s\n", err.Error())
					}
				}(ws)
			}
		}
	}(ws)
}

func (s *Server) HandleSocket(ws *websocket.Conn) {
	fmt.Printf("new incoming connection from %s\n", ws.RemoteAddr())
	s.conns[ws] = true
	func(ws *websocket.Conn) {
		buff := make([]byte, 1024)
		for {
			n, err := ws.Read(buff)
			if err != nil {
				if err == io.EOF {
					fmt.Printf("error end of file, %s\n", err.Error())
					break
				}

				fmt.Printf("read error, %s\n", err.Error())
				continue
			}

			msg := buff[:n]
			if string(msg) == "ping" {
				ws.Write([]byte("pong"))
			} else {
				ws.Write(msg)
			}
		}
	}(ws)
}
