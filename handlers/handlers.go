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
var repo = quotes.NewRepository(database.NewDb())

// Random takes care of the 'random' route.
func Random(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	err := writeJSON(w, parsedQuotes.Random())
	if err != nil {
		log.Fatal(err)
	}
}

// All takes care of the 'all' route.
func All(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	// parsedQuotes.Index(db)
	err := writeJSON(w, repo.All())
	if err != nil {
		log.Fatal(err)
	}
}

// Search takes care of the /search/{word} route.
func Search(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	word := strings.ToLower(p.ByName("word"))
	qs := parsedQuotes.List(word)

	if len(qs) != 0 {
		err := writeJSON(w, qs)
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
