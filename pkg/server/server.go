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

func WithRouter(router http.Handler) Option {
	return func(s *Server) {
		s.Router = router
	}
}

func NewServer(opts ...Option) (*Server, error) {
	var err error

	s := &Server{
		Host: "",
		Port: "",
	}

	for _, opt := range opts {
		opt(s)
	}

	s.Listener, err = net.Listen("tcp", s.Host+":"+s.Port)
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %w", err)
	}

	return s, nil
}

func (s *Server) Run() {
	fmt.Println("Starting server " + s.Host + ":" + s.Port)

	if err := http.Serve(s.Listener, s.Router); err != nil {
		fmt.Printf("Failed to serve: %v", err)
	}
}
