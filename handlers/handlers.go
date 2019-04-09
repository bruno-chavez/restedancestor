// Package handlers is used to separate the handlers from other functions of the API.
package handlers

import (
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/bruno-chavez/restedancestor/database"
	"github.com/bruno-chavez/restedancestor/quotes"
	"github.com/julienschmidt/httprouter"
)

const nbTop = 5

// parsedQuotes is a global variable to avoid multiple calls to Parser since always returns the same parsedQuotes.
var db = &database.File{}
var parsedQuotes = quotes.Parser(*db)

// Random takes care of the 'random' route.
func Random(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	err := writeJSON(w, parsedQuotes.Random())
	if err != nil {
		log.Fatal(err)
	}
}

// All takes care of the 'all' route.
func All(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	parsedQuotes.Index(db)
	err := writeJSON(w, parsedQuotes)
	if err != nil {
		log.Fatal(err)
	}
}

// Search takes care of the /search/{word} route.
func Search(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	word := p.ByName("word")
	matched := false
	quotesMatched := make([]quotes.QuoteType, 0)

	// quote is a QuoteType type.
	// xSlice is a parsedQuotes with one word and white-spaces that appear after a filter is applied.
	// xFilter is a string inside a wordXSlice.
	// This for is a candidate for been transformed into an recursive function inside lib.
	for _, quote := range parsedQuotes.Data {
		firstSlice := strings.Split(quote.Quote, " ")

		for _, firstFilter := range firstSlice {
			secondSlice := strings.Split(firstFilter, ",")

			for _, secondFilter := range secondSlice {
				thirdSlice := strings.Split(secondFilter, "!")

				for _, thirdFilter := range thirdSlice {
					filteredWord := strings.Split(thirdFilter, ".")

					// After all the filters are applied the filtered word is compared with the word that is been searched.
					// If a match is found, the quote that the filtered word belongs to, is printed.
					if filteredWord[0] == word {
						quotesMatched = append(quotesMatched, quote)
						matched = true
					}
				}
			}
		}
	}

	if matched {
		err := writeJSON(w, quotesMatched)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := writeNotFound(w, word)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Senile takes care of the 'senile' route
func Senile(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	data := parsedQuotes.Random()
	data1 := parsedQuotes.Random()

	quoteArray := strings.Split(data.Quote, " ")
	quoteArray1 := strings.Split(data1.Quote, " ")
	var quote string

	if len(quoteArray) < len(quoteArray1) {
		quote = stringModifier(quoteArray, quoteArray1)
	} else {
		quote = stringModifier(quoteArray1, quoteArray)
	}

	joinedQuote := quotes.QuoteType{Quote: quote}
	err := writeJSON(w, joinedQuote)
	if err != nil {
		log.Fatal(err)
	}
}

// Find takes care of the /one/{UUID} route.
func Find(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	uuidToSearch := p.ByName("uuid")

	offset, err := parsedQuotes.OffsetQuoteFromUUID(uuidToSearch)
	if err != nil {
		err = writeNotFound(w, uuidToSearch)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	quotes := parsedQuotes.Data
	err = writeJSON(w, quotes[*offset])
	if err != nil {
		log.Fatal(err)
	}
}

// Like takes care of the /one/{UUID}/like route.
func Like(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	uuidToSearch := p.ByName("uuid")

	if _, err := parsedQuotes.OffsetQuoteFromUUID(uuidToSearch); err != nil {
		err = writeNotFound(w, uuidToSearch)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	parsedQuotes.LikeQuote(db, uuidToSearch)

}

// Dislike takes care of the /one/{UUID}/dislike route.
func Dislike(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	uuidToSearch := p.ByName("uuid")

	if _, err := parsedQuotes.OffsetQuoteFromUUID(uuidToSearch); err != nil {
		err = writeNotFound(w, uuidToSearch)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	parsedQuotes.DislikeQuote(db, uuidToSearch)

}

// Top takes care of the /top route.
func Top(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	sort.Sort(parsedQuotes)

	i := 0
	top := make([]quotes.QuoteType, 0)
	for _, quote := range parsedQuotes.Data {
		if i >= nbTop {
			break
		}
		top = append(top, quote)
		i++
	}

	err := writeJSON(w, top)
	if err != nil {
		log.Fatal(err)
	}
}
