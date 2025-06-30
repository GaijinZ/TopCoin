package server

import (
	"fmt"
	"net"
	"net/http"
)

type Server struct {
	Listener net.Listener
	Router   http.Handler
	Host     string
	Port     string
}

type Option func(*Server)

func WithHost(host string) Option {
	return func(s *Server) {
		s.Host = host
	}
}

func WithPort(port string) Option {
	return func(s *Server) {
		s.Port = port
	}
}

func WithListener(listener net.Listener) Option {
	return func(s *Server) {
		s.Listener = listener
	}
}

func WithRouter(router http.Handler) Option {
	return func(s *Server) {
		s.Router = router
	}
}

func NewServer(opts ...Option) *Server {
	s := &Server{
		Host: "localhost",
		Port: "8080",
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

//func NewServer(listener net.Listener, router http.Handler) *Server {
//	return &Server{
//		Listener: listener,
//		Router:   router,
//	}
//}

func (s *Server) Run() {
	fmt.Println("Starting server")

	if err := http.Serve(s.Listener, s.Router); err != nil {
		fmt.Printf("Failed to serve: %v", err)
	}
}
