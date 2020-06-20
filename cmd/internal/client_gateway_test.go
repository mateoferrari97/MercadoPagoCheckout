package internal

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
)

type ClientStub struct {
	resp *http.Response
	err error
}

func (c *ClientStub) Do(_ *http.Request) (*http.Response, error) {
	if c.err != nil {
		return &http.Response{}, c.err
	}

	return c.resp, nil
}

func TestGateway_GetAccessToken(t *testing.T) {
	// Given
	c := &ClientStub{}
	g := &Gateway{Client: c}
	c.resp = &http.Response{
		Status:           "200",
		StatusCode:       200,
		Body:             ioutil.NopCloser(bytes.NewReader([]byte(`{"access_token": "1234"}`))),
	}
	// When
	accessToken, err := g.GetAccessToken(Credentials{
		ClientID:     "ABC123",
		ClientSecret: "123ABC",
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, accessToken, "1234")
}

func TestGateway_GetAccessToken_MercadoPagoError(t *testing.T) {
	// Given
	c := &ClientStub{}
	g := &Gateway{Client: c}
	c.resp = &http.Response{
		Status:           "500",
		StatusCode:       500,
		Body:             ioutil.NopCloser(bytes.NewReader([]byte(`{"error": "internal server error"}`))),
	}
	// When
	_, err := g.GetAccessToken(Credentials{
		ClientID:     "ABC123",
		ClientSecret: "123ABC",
	})

	// Then
	require.Error(t, err)
	require.EqualError(t, err, "{\"error\": \"internal server error\"}")
}

func TestGateway_GetAccessToken_UnmarshalError(t *testing.T) {
	// Given
	c := &ClientStub{}
	g := &Gateway{Client: c}
	c.resp = &http.Response{
		Status:           "200",
		StatusCode:       200,
		Body:             ioutil.NopCloser(bytes.NewReader([]byte(`{"access_token": 1123}`))),
	}
	// When
	_, err := g.GetAccessToken(Credentials{
		ClientID:     "ABC123",
		ClientSecret: "123ABC",
	})

	// Then
	require.Error(t, err)
	require.EqualError(t, err, "json: cannot unmarshal number into Go struct field .access_token of type string")
}

func TestGateway_CreatePreference(t *testing.T) {
	// Given
	c := &ClientStub{}
	g := &Gateway{Client: c}
	c.resp = &http.Response{
		Status:           "200",
		StatusCode:       200,
		Body:             ioutil.NopCloser(bytes.NewReader([]byte(`{"init_point": "https://mercadopago.com/checkout"}`))),
	}
	// When
	checkout, err := g.CreatePreference(newPreference())

	// Then
	require.NoError(t, err)
	require.Equal(t, checkout, "https://mercadopago.com/checkout")
}

func TestGateway_CreatePreference_MercadoPagoError(t *testing.T) {
	// Given
	c := &ClientStub{}
	g := &Gateway{Client: c}
	c.resp = &http.Response{
		Status:           "500",
		StatusCode:       500,
		Body:             ioutil.NopCloser(bytes.NewReader([]byte(`{"error": "internal server error"}`))),
	}
	// When
	_, err := g.CreatePreference(newPreference())

	// Then
	require.Error(t, err)
	require.EqualError(t, err, "{\"error\": \"internal server error\"}")
}

func TestGateway_CreatePreference_UnmarshalError(t *testing.T) {
	// Given
	c := &ClientStub{}
	g := &Gateway{Client: c}
	c.resp = &http.Response{
		Status:           "200",
		StatusCode:       200,
		Body:             ioutil.NopCloser(bytes.NewReader([]byte(`{"init_point": 1234}`))),
	}
	// When
	_, err := g.CreatePreference(newPreference())

	// Then
	require.Error(t, err)
	require.EqualError(t, err, "json: cannot unmarshal number into Go struct field .init_point of type string")
}

func newPreference() NewPreference {
	return NewPreference{
		Items: []Item{
			{
				Title:       "sherlock",
				Description: "holes",
				PictureURL:  "",
				Quantity:    1,
				UnitPrice:   15.75,
			},
		},
		Payer: Payer{
			Name:      "mateo",
			Surname:   "fc",
			Email:     "m@gmail.com",
			Phone:     Phone{
				AreaCode: "",
				Number:   "12345",
			},
			Address:   Address{
				ZipCode: "",
				Street:  "pepe",
				Number:  1234,
			},
			CreatedAt: "",
		},
		Redirect:   Redirect{
			Success: "http://baseurl.com/success",
			Pending: "http://baseurl.com/pending",
			failure: "http://baseurl.com/failure",
		},
		AutoReturn: true,
	}
}