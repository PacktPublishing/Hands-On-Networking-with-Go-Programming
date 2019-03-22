package router

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func getHandler(status int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write([]byte(strconv.Itoa(status)))
	})
}

func TestRouter(t *testing.T) {
	tests := []struct {
		name           string
		rtr            *Router
		r              *http.Request
		expectedStatus int
	}{
		{
			name:           "expect error",
			rtr:            New().AddRoute("/error", "GET", getHandler(500)).AddRoute("/success", "GET", getHandler(200)),
			r:              httptest.NewRequest("GET", "/error", nil),
			expectedStatus: 500,
		},
		{
			name:           "expect success",
			rtr:            New().AddRoute("/error", "GET", getHandler(500)).AddRoute("/success", "GET", getHandler(200)),
			r:              httptest.NewRequest("GET", "/success", nil),
			expectedStatus: 200,
		},
		{
			name:           "expect not found",
			rtr:            New().AddRoute("/error", "GET", getHandler(500)).AddRoute("/success", "GET", getHandler(200)),
			r:              httptest.NewRequest("GET", "/something", nil),
			expectedStatus: 404,
		},
		{
			name:           "method is important",
			rtr:            New().AddRoute("/error", "GET", getHandler(500)).AddRoute("/success", "GET", getHandler(200)),
			r:              httptest.NewRequest("POST", "/success", nil),
			expectedStatus: 404,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			test.rtr.ServeHTTP(w, test.r)
			if w.Result().StatusCode != test.expectedStatus {
				t.Errorf("expected status %d, got %d", test.expectedStatus, w.Result().StatusCode)
			}
		})
	}
}
