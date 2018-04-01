//Package lib contains various functions that may be useful for more than one part of the API.
package lib

import (
	"math/rand"
	"restedancestor/database"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

//Random returns a random quote from a QuoteList slice.
func Random(q []database.QuotesList) string {
	return q[rand.Intn(len(q))].Quotes
}
