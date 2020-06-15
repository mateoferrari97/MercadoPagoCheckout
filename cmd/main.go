package main

import (
	"github.com/mateoferrari97/mercadopago/cmd/internal"
	"github.com/mateoferrari97/mercadopago/cmd/server"
	"net/http"
)

const _defaultPort = ":8081"

func main() {
	server := server.NewServer()
	gateway := internal.NewClientGateway(&http.Client{})
	service := internal.NewController(gateway)
	handler := internal.NewHandler(service)

	server.HandleFunc("/ping", "GET", handler.Ping)
	server.HandleFunc("/access_token", "GET", handler.GetAccessToken)
	server.HandleFunc("/preferences", "POST", handler.CreatePreference)

	server.Run(_defaultPort)
}