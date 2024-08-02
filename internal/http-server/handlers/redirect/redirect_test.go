package redirect

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"rest_api_shortener/internal/http-server/handlers/redirect/mocks"
	"rest_api_shortener/internal/lib/api/response"
	"rest_api_shortener/internal/logger/handlers/slogdiscard"
	"testing"
)

func TestRedirectHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
		retCode   int
	}{
		{
			name:    "Success",
			alias:   "test_alias",
			url:     "https://google.com/",
			retCode: http.StatusFound,
		},
		{
			name:      "Empty alias",
			alias:     "",
			respError: "Invalid request alias",
			retCode:   http.StatusNotFound,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlGetMock := mocks.NewURLGetter(t)
			urlGetMock.On("GetURL", tc.alias).
				Return(tc.url, tc.mockError).
				Maybe()

			router := chi.NewRouter()
			router.Use(middleware.URLFormat)

			handler := New(slogdiscard.NewDiscardLogger(), urlGetMock)
			router.Get("/{alias}", handler)

			input := "" //fmt.Sprintf(`{"alias": "%s"}`, tc.alias)

			req, err := http.NewRequest(http.MethodGet, "/"+tc.alias, bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			retCode := http.StatusOK
			if tc.retCode != 0 {
				retCode = tc.retCode
			}
			require.Equal(t, rr.Code, retCode)

			body := rr.Body.String()

			var resp response.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
