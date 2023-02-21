package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"snippetbox/internal/assert"
	"testing"
)

func Test_secureHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	next := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("OK"))
	})

	secureHeaders(next).ServeHTTP(rr, r)

	rs := rr.Result()

	assert.Equal(t, rs.Header.Get("Content-Security-Policy"),
		"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
	assert.Equal(t, rs.Header.Get("Referrer-Policy"),
		"origin-when-cross-origin")
	assert.Equal(t, rs.Header.Get("X-Content-Type-Options"),
		"nosniff")
	assert.Equal(t, rs.Header.Get("X-Frame-Options"),
		"deny")
	assert.Equal(t, rs.Header.Get("X-XSS-Protection"),
		"0")
	assert.Equal(t, rs.StatusCode,
		http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, string(bytes.TrimSpace(body)), "OK")

	//
	//type args struct {
	//	next http.Handler
	//}
	//tests := []struct {
	//	name string
	//	args args
	//	want http.Handler
	//}{
	//	// TODO: Add test cases.
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		if got := secureHeaders(tt.args.next); !reflect.DeepEqual(got, tt.want) {
	//			t.Errorf("secureHeaders() = %v, want %v", got, tt.want)
	//		}
	//	})
	//}
}
