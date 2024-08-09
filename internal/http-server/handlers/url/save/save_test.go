package save

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"rest_api_shortener/internal/http-server/handlers/url/save/mocks"
	"rest_api_shortener/internal/logger/handlers/slogdiscard"
	"rest_api_shortener/internal/storage"
	"testing"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://google.com/",
		},
		{
			name:  "Empty alias",
			alias: "",
			url:   "https://google.com/",
		},
		{
			name:      "Empty URL",
			alias:     "test_alias",
			url:       "",
			respError: "Field URL is a required field",
		},
		{
			name:      "Invalid URL",
			alias:     "test_alias",
			url:       "test invalid URL",
			respError: "Field URL is not a valid URL",
		},
		{
			name:      "SaveURL Error",
			alias:     "test_alias",
			url:       "https://google.com/",
			respError: "Failed to save URL",
			mockError: errors.New("unexpected error"),
		},
		{
			name:      "Alias exist",
			alias:     "test_alias",
			url:       "https://google.com/",
			respError: "URL alias already exists",
			mockError: storage.ErrURLExists,
		},
		{
			name:      "Invalid alias JSON",
			alias:     "test_alias\"",
			url:       "https://google.com/",
			respError: "Filed to decode request",
		},
		{
			name:      "Invalid url JSON",
			alias:     "test_alias",
			url:       "https://google.com/\n",
			respError: "Filed to decode request",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlSaveMock := mocks.NewURLSaver(t)

			urlSaveMock.On("SaveURL", tc.url, mock.AnythingOfType("string")).
				Return(int64(1), tc.mockError).Maybe()

			handler := New(slogdiscard.NewDiscardLogger(), urlSaveMock)

			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)

			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp.Response))

			require.Equal(t, tc.respError, resp.Response.Error)
		})
	}
}
