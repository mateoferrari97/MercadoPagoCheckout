package internal

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ServiceStub struct {
	accessToken string
	checkout string
	err error
}

func (s *ServiceStub) GetAccessToken(_ string, _ string) (string, error) {
	return s.accessToken, s.err
}

func (s *ServiceStub) CreatePreference(_ string, _ NewPreference) (string, error) {
	return s.checkout, s.err
}

func TestHandler_GetAccessToken(t *testing.T) {
	// Given
	h := NewHandler(&ServiceStub{
		accessToken: "MY_ACCESS_TOKEN",
	})
	ts := httptest.NewServer(http.HandlerFunc(h.GetAccessToken))
	defer ts.Close()

	// When
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/access_token?client_id=MY_CLIENT_ID&client_secret=MY_CLIENT_SECRET", ts.URL), nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Then
	require.Equal(t, "MY_ACCESS_TOKEN", string(b))
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHandler_GetAccessToken_BadRequest_Error(t *testing.T) {
	tt := []struct{
		name string
		clientID string
		clientSecret string
		wantError string
	}{
		{
			name: "missing client id",
			clientSecret: "MY_CLIENT_SECRET",
			wantError: "client id is required",
		},
		{
			name: "missing client secret",
			clientID: "MY_CLIENT_ID",
			wantError: "client secret is required",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			h := NewHandler(&ServiceStub{
				accessToken: "MY_ACCESS_TOKEN",
			})
			ts := httptest.NewServer(http.HandlerFunc(h.GetAccessToken))
			defer ts.Close()

			// When
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/access_token?client_id=%s&client_secret=%s", ts.URL, tc.clientID, tc.clientSecret), nil)
			if err != nil {
				t.Fatal(err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			errorMessage, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			// Then
			require.Equal(t, tc.wantError, string(errorMessage))
			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func TestHandler_GetAccessToken_Error(t *testing.T) {
	tt := []struct{
		name string
		err error
		wantError string
		wantErrorStatusCode int
	}{
		{
			name: "bad request from server",
			err: NewError("bad request", http.StatusBadRequest),
			wantError: "couldn't get access token: bad request",
			wantErrorStatusCode: http.StatusBadRequest,
		},
		{
			name: "unauthorized client from server",
			err: NewError("unauthorized", http.StatusUnauthorized),
			wantError: "couldn't get access token: unauthorized",
			wantErrorStatusCode: http.StatusUnauthorized,
		},
		{
			name: "internal server error from server",
			err: NewError("internal server error", http.StatusInternalServerError),
			wantError: "couldn't get access token: internal server error",
			wantErrorStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			h := NewHandler(&ServiceStub{
				err: tc.err,
			})
			ts := httptest.NewServer(http.HandlerFunc(h.GetAccessToken))
			defer ts.Close()

			// When
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/access_token?client_id=MY_CLIENT_ID&client_secret=MY_CLIENT_SECRET", ts.URL), nil)
			if err != nil {
				t.Fatal(err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			errorMessage, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			// Then
			require.Equal(t, tc.wantError, string(errorMessage))
			require.Equal(t, tc.wantErrorStatusCode, resp.StatusCode)
		})
	}
}

func TestHandler_CreatePreference(t *testing.T) {
	// Given
	h := NewHandler(&ServiceStub{
		checkout: "https://mercadopago.com/MY_CHECKOUT_PATH",
	})
	body := []byte(`{
		"items": [
			{
				"title": "Libro Sherlock Holmes 1era edicion",
				"description": "Nuevo libro de sherlock holmes 2020",
				"quantity": 1,
				"unit_price": 150.70,
				"picture_url": "https://www.comunidadbaratz.com/wp-content/uploads/Instrucciones-a-tener-en-cuenta-sobre-como-se-abre-un-libro-nuevo.jpg"
			}
   	 	],
		"payer": {
			"name": "Mateo",
			"email": "mateo.ferrari@gmail.com",
			"phone": {
				"number": "11111111"
			},
			"address": {
				"street": "posta",
				"number": 4789
			},
			"date_created": "14-06-2020"
		}
	}`)
	ts := httptest.NewServer(http.HandlerFunc(h.CreatePreference))
	defer ts.Close()

	// When
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/preferences", ts.URL), bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("access_token", "MY_ACCESS_TOKEN")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Then
	require.Equal(t, "https://mercadopago.com/MY_CHECKOUT_PATH", string(b))
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHandler_CreatePreference_UnprocessableEntity_Error(t *testing.T) {
	// Given
	h := NewHandler(&ServiceStub{
		checkout: "https://mercadopago.com/MY_CHECKOUT_PATH",
	})
	body := []byte(`{
		"items": [
			{
				"title": "Libro Sherlock Holmes 1era edicion",
				"description": "Nuevo libro de sherlock holmes 2020",
				"quantity": "1",
				"unit_price": 150.70,
				"picture_url": "https://www.comunidadbaratz.com/wp-content/uploads/Instrucciones-a-tener-en-cuenta-sobre-como-se-abre-un-libro-nuevo.jpg"
			}
   	 	],
		"payer": {
			"name": "Mateo",
			"email": "mateo.ferrari@gmail.com",
			"phone": {
				"number": "11111111"
			},
			"address": {
				"street": "posta",
				"number": 4789
			},
			"date_created": "14-06-2020"
		}
	}`)
	ts := httptest.NewServer(http.HandlerFunc(h.CreatePreference))
	defer ts.Close()

	// When
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/preferences", ts.URL), bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("access_token", "MY_ACCESS_TOKEN")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Then
	require.Equal(t, "couldn't decode body: json: cannot unmarshal string into Go struct field Item.items.quantity of type int", string(b))
	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}

func TestHandler_CreatePreference_BadRequest_Error(t *testing.T) {
	tt := []struct{
		name string
		body []byte
		wantError string
	}{
		{
			name: "missing items field",
			body: []byte(`{
					"payer": {
						"name": "Mateo",
						"email": "mateo.ferrari@gmail.com",
						"phone": {
							"number": "11111111"
						},
						"address": {
							"street": "posta",
							"number": 4789
						},
						"date_created": "14-06-2020"
					}
			}`),
			wantError: "validation error: Key: 'NewPreference.Items' Error:Field validation for 'Items' failed on the 'required' tag",
		},
		{
			name: "items with length 0",
			body: []byte(`{
					"items": [],
					"payer": {
						"name": "Mateo",
						"email": "mateo.ferrari@gmail.com",
						"phone": {
							"number": "11111111"
						},
						"address": {
							"street": "posta",
							"number": 4789
						},
						"date_created": "14-06-2020"
					}
			}`),
			wantError: "validation error: Key: 'NewPreference.Items' Error:Field validation for 'Items' failed on the 'min' tag",
		},
		{
			name: "no payer sent",
			body: []byte(`{
					"items": [        
								{
									"title": "Libro Sherlock Holmes 1era edicion",
									"description": "Nuevo libro de sherlock holmes 2020",
									"quantity": 1,
									"unit_price": 150.70,
									"picture_url": "https://www.comunidadbaratz.com/wp-content/uploads/Instrucciones-a-tener-en-cuenta-sobre-como-se-abre-un-libro-nuevo.jpg"
								}
					]
			}`),
			wantError: "validation error: Key: 'NewPreference.Payer.Name' Error:Field validation for 'Name' failed on the 'required' tag\nKey: 'NewPreference.Payer.Email' Error:Field validation for 'Email' failed on the 'required' tag\nKey: 'NewPreference.Payer.Phone.Number' Error:Field validation for 'Number' failed on the 'required' tag\nKey: 'NewPreference.Payer.Address.Street' Error:Field validation for 'Street' failed on the 'required' tag\nKey: 'NewPreference.Payer.Address.Number' Error:Field validation for 'Number' failed on the 'required' tag\nKey: 'NewPreference.Payer.CreatedAt' Error:Field validation for 'CreatedAt' failed on the 'required' tag",
		},
		{
			name: "missing name inside payer field",
			body: []byte(`{
					"items": [
						{
							"title": "Libro Sherlock Holmes 1era edicion",
							"description": "Nuevo libro de sherlock holmes 2020",
							"quantity": 1,
							"unit_price": 150.70,
							"picture_url": "https://www.comunidadbaratz.com/wp-content/uploads/Instrucciones-a-tener-en-cuenta-sobre-como-se-abre-un-libro-nuevo.jpg"
						}
					],
					"payer": {
						"email": "mateo.ferrari@gmail.com",
						"phone": {
							"number": "11111111"
						},
						"address": {
							"street": "posta",
							"number": 4789
						},
						"date_created": "14-06-2020"
					}
			}`),
			wantError: "validation error: Key: 'NewPreference.Payer.Name' Error:Field validation for 'Name' failed on the 'required' tag",
		},
		{
			name: "missing unit_price inside items field",
			body: []byte(`{
					"items": [
						{
							"title": "Libro Sherlock Holmes 1era edicion",
							"description": "Nuevo libro de sherlock holmes 2020",
							"quantity": 1,
							"picture_url": "https://www.comunidadbaratz.com/wp-content/uploads/Instrucciones-a-tener-en-cuenta-sobre-como-se-abre-un-libro-nuevo.jpg"
						}
					],
					"payer": {
						"name": "Mateo",
						"email": "mateo.ferrari@gmail.com",
						"phone": {
							"number": "11111111"
						},
						"address": {
							"street": "posta",
							"number": 4789
						},
						"date_created": "14-06-2020"
					}
			}`),
			wantError: "validation error: Key: 'Item.UnitPrice' Error:Field validation for 'UnitPrice' failed on the 'required' tag",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			h := NewHandler(&ServiceStub{})
			ts := httptest.NewServer(http.HandlerFunc(h.CreatePreference))
			defer ts.Close()

			// When
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/preferences", ts.URL), bytes.NewReader(tc.body))
			if err != nil {
				t.Fatal(err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			// Then
			require.Equal(t, tc.wantError, string(b))
			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func TestHandler_CreatePreference_Unauthorized_Error(t *testing.T) {
	// Given
	h := NewHandler(&ServiceStub{})
	ts := httptest.NewServer(http.HandlerFunc(h.CreatePreference))
	defer ts.Close()

	body := []byte(`{
		"items": [
			{
				"title": "Libro Sherlock Holmes 1era edicion",
				"description": "Nuevo libro de sherlock holmes 2020",
				"quantity": 1,
				"unit_price": 150.70,
				"picture_url": "https://www.comunidadbaratz.com/wp-content/uploads/Instrucciones-a-tener-en-cuenta-sobre-como-se-abre-un-libro-nuevo.jpg"
			}
   	 	],
		"payer": {
			"name": "Mateo",
			"email": "mateo.ferrari@gmail.com",
			"phone": {
				"number": "11111111"
			},
			"address": {
				"street": "posta",
				"number": 4789
			},
			"date_created": "14-06-2020"
		}
	}`)

	// When
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/preferences", ts.URL), bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Then
	require.Equal(t, "access token is required", string(b))
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestHandler_CreatePreference_ClientError(t *testing.T) {
	tt := []struct{
		name string
		err error
		wantError string
		wantErrorStatusCode int
	}{
		{
			name: "bad request from server",
			err: NewError("bad request", http.StatusBadRequest),
			wantError: "couldn't create checkout: bad request",
			wantErrorStatusCode: http.StatusBadRequest,
		},
		{
			name: "unauthorized client from server",
			err: NewError("unauthorized", http.StatusUnauthorized),
			wantError: "couldn't create checkout: unauthorized",
			wantErrorStatusCode: http.StatusUnauthorized,
		},
		{
			name: "internal server error from server",
			err: NewError("internal server error", http.StatusInternalServerError),
			wantError: "couldn't create checkout: internal server error",
			wantErrorStatusCode: http.StatusInternalServerError,
		},
		{
			name: "couldn't cast error",
			err: errors.New("random error"),
			wantError: "couldn't create checkout: random error",
			wantErrorStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			h := NewHandler(&ServiceStub{
				err: tc.err,
			})
			body := []byte(`{
					"items": [
						{
							"title": "Libro Sherlock Holmes 1era edicion",
							"description": "Nuevo libro de sherlock holmes 2020",
							"quantity": 1,
							"unit_price": 150.70,
							"picture_url": "https://www.comunidadbaratz.com/wp-content/uploads/Instrucciones-a-tener-en-cuenta-sobre-como-se-abre-un-libro-nuevo.jpg"
						}
					],
					"payer": {
						"name": "Mateo",
						"email": "mateo.ferrari@gmail.com",
						"phone": {
							"number": "11111111"
						},
						"address": {
							"street": "posta",
							"number": 4789
						},
						"date_created": "14-06-2020"
					}
				}`)
			ts := httptest.NewServer(http.HandlerFunc(h.CreatePreference))
			defer ts.Close()

			// When
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/preferences", ts.URL), bytes.NewReader(body))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Add("access_token", "MY_ACCESS_TOKEN")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			// Then
			require.Equal(t, tc.wantError, string(b))
			require.Equal(t, tc.wantErrorStatusCode, resp.StatusCode)
		})
	}
}
