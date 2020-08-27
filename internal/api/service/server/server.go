package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Server ...
type Server struct {
	server *http.Server
}

// New ...
func New(handler http.Handler, timeout time.Duration, port int) (*Server, error) {
	switch {
	case handler == nil:
		return nil, errors.New("empty handler")
	case timeout < 0:
		return nil, fmt.Errorf("invalid timeout: %d", timeout)
	case port < 1 || port > 65535:
		return nil, fmt.Errorf("invalid port: %d", port)
	}

	return &Server{
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			Handler:      handler,
			ReadTimeout:  timeout,
			WriteTimeout: timeout,
		},
	}, nil
}

// Start ...
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}
