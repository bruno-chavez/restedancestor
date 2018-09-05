package database_test

import (
	"testing"

	"github.com/bruno-chavez/restedancestor/database"
)

func TestRandom(t *testing.T) {
	a := make(database.QuoteSlice, 0)
	quote := database.QuoteType{
		Quote: "this is a test",
	}
	a = append(a, quote)
	rand := a.Random()
	if a.Random() != quote {
		t.Errorf("Unexpected value %s, expected %s", rand, quote)
	}
}
