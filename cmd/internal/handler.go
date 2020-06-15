package internal

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
)

var _v = validator.New()

type Service interface {
	GetAccessToken(clientID string, clientSecret string) (string, error)
	CreatePreference(preference NewPreference) (string, error)
}

type Handler struct {
	Service Service
}

func NewHandler(service Service) *Handler{
	return &Handler{
		Service: service,
	}
}

func (h *Handler) Ping(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "pong")
}

func (h *Handler) GetAccessToken(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	if clientID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "client id is required")
		return
	}

	clientSecret := r.URL.Query().Get("client_secret")
	if clientSecret == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "client secret is required")
		return
	}

	accessToken, err := h.Service.GetAccessToken(clientID, clientSecret)
	if err != nil {
		w.WriteHeader(getStatusCodeFromError(err))
		fmt.Fprintf(w, "couldn't get access token: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", accessToken)
}

func (h *Handler) CreatePreference(w http.ResponseWriter, r *http.Request) {
	var preference NewPreference
	if err := json.NewDecoder(r.Body).Decode(&preference); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "couldn't decode body: %v", err)
		return
	}

	if err := _v.Struct(preference); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, fmt.Sprintf("validation error: %v", err))
	}

	checkoutURL, err := h.Service.CreatePreference(preference)
	if err != nil {
		w.WriteHeader(getStatusCodeFromError(err))
		fmt.Fprintf(w, "couldn't create checkout: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, fmt.Sprintf("%s", checkoutURL))
}

func getStatusCodeFromError(err error) int {
	e, ok := err.(*Error)
	if !ok {
		return http.StatusInternalServerError
	}
	
	return e.StatusCode
}



