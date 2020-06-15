package internal

import (
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

func (s *ServiceStub) CreatePreference(_ NewPreference) (string, error) {
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
