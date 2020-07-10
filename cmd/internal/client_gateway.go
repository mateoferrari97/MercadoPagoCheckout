package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const _baseURL = "https://api.mercadopago.com"

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type Gateway struct {
	Client Client
}

func NewClientGateway(client Client) *Gateway {
	return &Gateway{
		Client: client,
	}
}

func (g *Gateway) GetAccessToken(credentials Credentials) (string, error) {
	path := &url.Values{}
	path.Add("client_id", credentials.ClientID)
	path.Add("client_secret", credentials.ClientSecret)
	path.Add("grant_type", "client_credentials")
	queryParams := path.Encode()

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s%s", _baseURL, "/oauth/token?", queryParams), nil)
	if err != nil {
		return "", err
	}

	resp, err := g.Client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		return "", NewError(string(body), resp.StatusCode)
	}

	var r struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.Unmarshal(body, &r); err != nil {
		return "", err
	}

	return r.AccessToken, nil
}

func (g *Gateway) CreatePreference(accessToken string, preference NewPreference) (string, error) {
	queryValues := &url.Values{}
	queryValues.Add("access_token", accessToken)
	queryParams := queryValues.Encode()

	b, err := json.Marshal(preference)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s%s", _baseURL, "/checkout/preferences?", queryParams), bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	resp, err := g.Client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		return "", NewError(string(body), resp.StatusCode)
	}

	var r struct {
		CheckoutURL string `json:"init_point"`
	}

	if err := json.Unmarshal(body, &r); err != nil {
		return "", err
	}

	return r.CheckoutURL, nil
}
