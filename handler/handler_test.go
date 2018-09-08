package handler_test

import (
	"net/http"
	"testing"

	h "github.com/bruno-chavez/restedancestor/handler"
)

var writed []byte

type StubWriter struct {
}

func (f StubWriter) Header() http.Header {
	h := http.Header{}

	return h
}

func (f StubWriter) Write(b []byte) (int, error) {
	writed = b

	return 0, nil
}

func (f StubWriter) WriteHeader(statusCode int) {

}

func TestRandom(t *testing.T) {

	w := StubWriter{}
	h.AHandler(w)
	if len(writed) == 0 {
		t.Errorf("Unexpected value %s, expected %s", "a", "z")
	}
}
