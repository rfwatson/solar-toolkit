package handler_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"git.netflux.io/rob/solar-toolkit/gateway/handler"
	"git.netflux.io/rob/solar-toolkit/inverter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type store struct {
	err error
}

func (s *store) InsertETRuntimeData(*inverter.ETRuntimeData) error {
	return s.err
}

func TestHandler(t *testing.T) {
	testCases := []struct {
		name           string
		httpMethod     string
		path           string
		body           string
		storeErr       error
		wantStatusCode int
		wantBody       string
	}{
		{
			name:           "method not allowed",
			httpMethod:     http.MethodGet,
			path:           "/gateway/et_runtime_data",
			wantStatusCode: http.StatusMethodNotAllowed,
			wantBody:       "method not allowed\n",
		},
		{
			name:           "method not allowed",
			httpMethod:     http.MethodPost,
			path:           "/gateway/foo",
			wantStatusCode: http.StatusNotFound,
			wantBody:       "endpoint not found\n",
		},
		{
			name:           "invalid payload",
			httpMethod:     http.MethodPost,
			path:           "/gateway/et_runtime_data",
			body:           `{`,
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       "unexpected error\n",
		},
		{
			name:           "invalid timestamp",
			httpMethod:     http.MethodPost,
			path:           "/gateway/et_runtime_data",
			body:           `{"timestamp": "1970-01-01T00:00:00Z"}`,
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "invalid data\n",
		},
		{
			name:           "store error",
			httpMethod:     http.MethodPost,
			path:           "/gateway/et_runtime_data",
			body:           `{"timestamp": "2022-01-01T00:00:00Z"}`,
			storeErr:       errors.New("boom"),
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       "unexpected error\n",
		},
		{
			name:           "OK",
			httpMethod:     http.MethodPost,
			path:           "/gateway/et_runtime_data",
			body:           `{"timestamp": "2022-01-01T00:00:00Z"}`,
			wantStatusCode: http.StatusOK,
			wantBody:       "OK\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockStore := store{err: tc.storeErr}
			handler := handler.New(&mockStore)
			req := httptest.NewRequest(tc.httpMethod, tc.path, strings.NewReader(tc.body))
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			resp := rec.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.wantStatusCode, resp.StatusCode)

			if tc.wantBody != "" {
				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				assert.Equal(t, tc.wantBody, string(body))
			}
		})
	}
}
