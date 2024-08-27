package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}
func (s *Server) handleWSOrderbook(ws *websocket.Conn){
	fmt.Println("new incomming connections",ws.RemoteAddr())
	for {
		payload: fmt.Sprintf("orderedbook data -> %d\n",time.Now().UnixNano())
		ws.Write([]byte(payload))
		time.Sleep(time.Second*2)
	}
}
func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("new incoming connection from client:", ws.RemoteAddr())
	s.conns[ws] = true
	s.readLoop(ws)
}
func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read error", err)
			continue
		}
		msg := buf[:n]
		s.broadcast(msg)
		fmt.Println(string(msg))
		ws.Write([]byte("thank you for the msg"))
	}
}
func (s *Server) broadcast(b []byte){
	for ws :=range s.conns{
		go func(){
			if _, err :=ws.Write(b);err!=nil{
				fmt.Println("write error",err)
			}
		}(ws)
	}
}
func main() {
server :=NewServer()
http.Handle('/ws',websocket.Handler(server.handleWS))
http.Handle("/orderedbookfeed",websocket.Hanlder(server.handleWSOrderbook()))
http.ListenAndServe(":3000",nil)
}
