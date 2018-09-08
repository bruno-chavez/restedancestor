package handler_test

import (
	"net/http"
	"testing"

	"github.com/bruno-chavez/restedancestor/handler"
)

var (
	writed []byte
	h      http.Header
)

type StubWriter struct {
}

func (f StubWriter) Header() http.Header {
	h = http.Header{}

	return h
}

func (f StubWriter) Write(b []byte) (int, error) {
	writed = b

	return 0, nil
}

func (f StubWriter) WriteHeader(statusCode int) {
}

func TestRandomOptions(t *testing.T) {
	handlerOptions(t, handler.RandomHandler)
}

func TestAllOptions(t *testing.T) {
	handlerOptions(t, handler.AllHandler)
}

func TestSearchOptions(t *testing.T) {
	handlerOptions(t, handler.SearchHandler)
}

func handlerOptions(t *testing.T, fn func(w http.ResponseWriter, r *http.Request)) {
	w := StubWriter{}
	r := http.Request{
		Method: "OPTIONS",
	}
	fn(w, &r)
	if len(h.Get("Allow")) == 0 {
		t.Error("Unexpected void value for « Allow » header")
	}
}
