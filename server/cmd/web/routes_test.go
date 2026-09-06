package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
)

func TestClientIPRateLimitKey(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   string
	}{
		{name: "IPv4", header: "203.0.113.42", want: "203.0.113.42"},
		{name: "IPv6", header: "2001:db8:1234:5678::1", want: "2001:db8:1234:5678::"},
		{name: "rightmost proxy value", header: "198.51.100.10, 203.0.113.42", want: "203.0.113.42"},
		{name: "missing header fails closed", want: ""},
		{name: "invalid header fails closed", header: "not-an-ip", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			handler := middleware.ClientIPFromXFF()(http.HandlerFunc(
				func(_ http.ResponseWriter, r *http.Request) {
					var err error
					got, err = clientIPRateLimitKey(r)
					if err != nil {
						t.Errorf("clientIPRateLimitKey() error = %v", err)
					}
				},
			))

			request := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.header != "" {
				request.Header.Set("X-Forwarded-For", tt.header)
			}
			handler.ServeHTTP(httptest.NewRecorder(), request)

			if got != tt.want {
				t.Fatalf("clientIPRateLimitKey() = %q, want %q", got, tt.want)
			}
		})
	}
}
