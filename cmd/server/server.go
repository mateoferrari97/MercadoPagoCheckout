package server

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Server struct {
	server *mux.Router
}

func NewServer() *Server{
	return &Server{server: mux.NewRouter()}
}

func (s *Server) Run(port string) {
	if port == "" {
		port = "8081"
		log.Printf("defaulting to port %s", port)
	}

	if string(port[0]) == ":" {
		port = port[1:]
	}

	log.Printf("Listening on port %s", port)
	http.ListenAndServe(":" + port, s.server)
}

func (s *Server) HandleFunc(path string, method string, h http.HandlerFunc) {
	s.server.HandleFunc(path, h).Methods(method)
}