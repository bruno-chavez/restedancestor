package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func stringModifier(minQuote []string, maxQuote []string) string {
	// Generate the number of words to be modified
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	swappable := r1.Intn(len(minQuote)-1) + 1

	for i := 0; i < swappable; i++ {
		si := rand.NewSource(time.Now().UnixNano() * 10000)
		ri := rand.New(si)

		// Choose a random place in each array and decide whether replacing or inserting
		minPlace := ri.Intn(len(minQuote)) * 7 % len(minQuote)
		maxPlace := ri.Intn(len(maxQuote)) * 9 % len(maxQuote)
		modifier := ri.Float64()

		// Replace or Insert
		if modifier > .5 {
			maxQuote[maxPlace] = minQuote[minPlace]
		} else {
			maxQuote = append(maxQuote, "")
			copy(maxQuote[maxPlace+1:], maxQuote[maxPlace:])
			maxQuote[maxPlace] = minQuote[minPlace]
		}
	}

	return strings.Join(maxQuote, " ")
}

func writeJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	marshaledData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	_, err = w.Write(marshaledData)
	if err != nil {
		return err
	}

	return nil
}

func writeNotFound(w http.ResponseWriter, message string) error {

	// Sets NotFound header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	// Writes NotFound message
	notFound := map[string]string{"code": "404", "message": "'" + message + "' was not found"}
	notFoundJSON, err := json.Marshal(notFound)
	if err != nil {
		return err
	}
	_, err = w.Write(notFoundJSON)
	if err != nil {
		return err
	}

	return nil
}

/*
//BadRequest returns a ready to use badrequest text, not used for now
func BadRequest(requestText string) []byte{
	badRequest := map[string]string{"code": "400", "message": requestText}
	badRequestJson, _ := json.MarshalIndent(badRequest, "", "")
	return badRequestJson
}
*/

// NotFound returns a ready to write message for ResponseWriter when needed a 404 error.

// Filter returns a slice with all the words of a QuoteSlice after been filtered.
/*func Filter(filters []string, slice database.QuoteSlice) []string {
	var count int
	filteredWords := make([]string, 0)

		for _, quote := range slice {
			s := strings.Split(quote.Quote, " ")
			for _, word := range s {
				filtered := word
				for count < len(filters) {
					filtered = strings.Split(filtered, filters[count])[0]
					count++
				}
				count = 0
				filteredWords = append(filteredWords, filtered)
			}
	}
	return filteredWords
}*/
