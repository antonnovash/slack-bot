package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"slack-bot/slack-bot/pkg/controller"
)

// TODO: add logger

// Server handles incoming requests from Google Forms.
type Server struct {
	server     *http.Server
	controller *controller.Controller
}

// New creates a new instance of Server which is HTTP server with custom handler and a controller.
func New(cfg Config, c *controller.Controller) (*Server, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config is invalid: %v", err)
	}
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("controller is invalid: %v", err)
	}
	s := &Server{}
	r := mux.NewRouter()
	//r.HandleFunc("/",s.Handler)
	r.HandleFunc("/", s.Handler).Methods(http.MethodPost)
	r.HandleFunc("/auth",s.Auth).Methods(http.MethodGet)
	//r.Methods(http.MethodPost).Path("/").HandlerFunc(s.Handler)
	s.controller = c
	log.Println(cfg.Address)
	s.server = &http.Server{
		Addr:    cfg.Address,
		Handler: r,
	}
	return s, nil
}

// Run starts an HTTP server under the hood of Server.
func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		err := s.server.Shutdown(ctx)
		if err != nil {
			log.Printf("could not gracefully shut down: %v", err)
		}
	}()
	log.Printf("Listen on %s", s.server.Addr)
	err := s.server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("could not listen: %v", err)
	}
	return nil
}
