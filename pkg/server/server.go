package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"slack-bot/slack-bot/pkg/controller"

	"github.com/gorilla/mux"
)

// TODO: add logger

// Server handles incoming requests from Google Forms.
type Server struct {
	server     *http.Server
	controller *controller.Controller
	chToken    chan<- string
}

// New creates a new instance of Server which is HTTP server with custom handler and a controller.
func New(Address string, chToken chan<- string) (*Server, error) {
	//TODO validation
	s := &Server{chToken: chToken}
	r := mux.NewRouter()
	//r.HandleFunc("/",s.Handler)
	r.HandleFunc("/", s.Handler).Methods(http.MethodPost)
	r.HandleFunc("/add", s.addToSlack).Methods(http.MethodConnect, http.MethodGet)
	r.HandleFunc("/auth", s.auth).Methods(http.MethodGet)
	r.HandleFunc("/home", s.home).Methods(http.MethodGet)
	log.Println(Address)
	s.server = &http.Server{
		Addr:    Address,
		Handler: r,
	}
	return s, nil
}

// Run starts an HTTP server under the hood of Server.
func (s *Server) Run(ctx context.Context) error {
	log.Println("server run")
	go func() {
		<-ctx.Done()
		err := s.server.Shutdown(context.Background())
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
