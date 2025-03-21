package api

import (
	"net/http"
	"net/http/httptest"
)

func mockServerWithoutBodyAndVerifier(status int, verifier func(request *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(status)
		verifier(request)
	}))
}

func mockServerWithBodyAndVerifier(status int, body []byte, verifier func(request *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(status)
		_, _ = writer.Write(body)
		verifier(request)
	}))
}
