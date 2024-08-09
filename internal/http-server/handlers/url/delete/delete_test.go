package delete

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"rest_api_shortener/internal/http-server/handlers/url/delete/mocks"
	"rest_api_shortener/internal/lib/api/response"
	"rest_api_shortener/internal/logger/handlers/slogdiscard"
	"rest_api_shortener/internal/storage"
	"testing"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
		},
		{
			name:      "Empty alias",
			alias:     "",
			respError: "Invalid request alias",
			mockError: errors.New("Invalid request alias"),
		},
		{
			name:      "Invalid alias JSON",
			alias:     "invalid JSON\"",
			respError: "Filed to decode request",
			mockError: errors.New("Filed to decode request"),
		},
		{
			name:      "Alias not found",
			alias:     "test_alias",
			respError: "URL alias not found",
			mockError: storage.ErrURLNotFound,
		},
		{
			name:      "Failed delete alias",
			alias:     "test_alias",
			respError: "Failed to delete URL alias",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlDeleteMock := mocks.NewURLDeleter(t)

			urlDeleteMock.On("DeleteURL", tc.alias).
				Return(tc.mockError).Maybe()

			handler := New(slogdiscard.NewDiscardLogger(), urlDeleteMock)

			input := fmt.Sprintf(`{"alias": "%s"}`, tc.alias)

			req, err := http.NewRequest(http.MethodDelete, "/{alias}", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp response.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
