package server

import (
	"fmt"
	"net"
	"net/http"
)

type Server struct {
	Listener net.Listener
	Router   http.Handler
}

func NewServer(listener net.Listener, router http.Handler) *Server {
	return &Server{
		Listener: listener,
		Router:   router,
	}
}

func (s *Server) Run() {
	fmt.Println("Starting server")

	if err := http.Serve(s.Listener, s.Router); err != nil {
		fmt.Printf("Failed to serve: %v", err)
	}
}
