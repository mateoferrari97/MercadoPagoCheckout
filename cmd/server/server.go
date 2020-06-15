package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	server *mux.Router
}

func NewServer() *Server{
	return &Server{server: mux.NewRouter()}
}

func (s *Server) Run(port string) {
	http.ListenAndServe(port, s.server)
}

func (s *Server) HandleFunc(path string, method string, h http.HandlerFunc) {
	s.server.HandleFunc(path, h).Methods(method)
}