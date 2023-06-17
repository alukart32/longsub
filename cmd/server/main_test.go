package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSubstring(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "Empty string",
			arg:  "",
			want: "",
		},
		{
			name: "Get b for bbbb",
			arg:  "bbbb",
			want: "b",
		},
		{
			name: "Get abc for abcabcbb",
			arg:  "abcabcbb",
			want: "abc",
		},
		{
			name: "Get wke for pwwkew",
			arg:  "pwwkew",
			want: "wke",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := getSubstring(tt.arg)
			assert.EqualValues(t, tt.want, actual)
		})
	}
}

func TestLongSubstringHandler(t *testing.T) {
	type req struct {
		body        string
		method      string
		contentType string
	}
	type want struct {
		response    string
		contentType string
		code        int
	}
	tests := []struct {
		want want
		req  req
		name string
	}{
		{
			name: "Http method GET, status code 405",
			req: req{
				method: http.MethodGet,
			},
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "Content-type application/json, status code 415",
			req: req{
				method:      http.MethodPost,
				contentType: "application/json",
			},
			want: want{
				code: http.StatusUnsupportedMediaType,
			},
		},
		{
			name: "Empty request, status code 400",
			req: req{
				body:        "",
				method:      http.MethodPost,
				contentType: "text/plain",
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "Respond b for bbbb",
			req: req{
				body:        "bbbb",
				method:      http.MethodPost,
				contentType: "text/plain",
			},
			want: want{
				code:        http.StatusOK,
				response:    "b",
				contentType: "text/plain",
			},
		},
		{
			name: "Respond abc for abcabcbb",
			req: req{
				body:        "abcabcbb",
				method:      http.MethodPost,
				contentType: "text/plain",
			},
			want: want{
				code:        http.StatusOK,
				response:    "abc",
				contentType: "text/plain",
			},
		},
		{
			name: "Respond wke for pwwkew",
			req: req{
				body:        "pwwkew",
				method:      http.MethodPost,
				contentType: "text/plain",
			},
			want: want{
				code:        http.StatusOK,
				response:    "wke",
				contentType: "text/plain",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.req.method,
				"/api/substring", bytes.NewBufferString(tt.req.body))
			req.Header.Set("Content-Type", tt.req.contentType)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(LongSubstringHandler)
			h.ServeHTTP(w, req)
			resp := w.Result()
			defer resp.Body.Close()

			assert.EqualValues(t, tt.want.code, resp.StatusCode,
				"expected status code %d, got %d", tt.want.code, w.Code)

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			if len(tt.want.contentType) != 0 {
				assert.EqualValues(t, tt.want.contentType, resp.Header.Get("Content-Type"),
					"expected content-type %s, got %s", tt.want.contentType, resp.Header.Get("Content-Type"))
			}
			if len(tt.want.response) != 0 {
				assert.EqualValues(t, tt.want.response, string(respBody),
					"expected body %s, got %s", tt.want.response, w.Body.String())
			}
		})
	}
}
