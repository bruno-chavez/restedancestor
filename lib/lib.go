//Package lib contains various functions that may be useful for more than one part of the API.
package lib

import (
	"github.com/bruno-chavez/restedancestor/database"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

//Random returns a random quote from a QuoteList slice.
func Random(q []database.QuotesList) database.QuotesList {
	return q[rand.Intn(len(q))]
}

/*
//BadRequest returns a ready to use badrequest text, not used for now
func BadRequest(requestText string) []byte{
	badRequest := map[string]string{"code": "400", "message": requestText}
	badRequestJson, _ := json.MarshalIndent(badRequest, "", "")
	return badRequestJson
}
*/