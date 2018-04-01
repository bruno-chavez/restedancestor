//Package database takes care of properly handle the database to be used in other parts of the API.
package database

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type QuotesList struct {
	Quotes string `json:"quote"`
}

//Parser fetches from database.json and puts it on a slice.
func Parser() []QuotesList {

	rawjson, _ := os.Open("database/database.json")
	readjson, _ := ioutil.ReadAll(rawjson)

	//313 is the total numer of quotes.
	parsedjson := make([]QuotesList, 313)
	json.Unmarshal(readjson, &parsedjson)

	return parsedjson
}
