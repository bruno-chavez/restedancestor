// Package handlers is used to separate the handlers from other functions of the API.
package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/bruno-chavez/restedancestor/database"
	"github.com/bruno-chavez/restedancestor/quotes"
	"github.com/julienschmidt/httprouter"
)

var repo = quotes.NewRepository(database.NewDb())

// Random takes care of the 'random' route.
func Random(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	q := repo.Random()
	err := writeJSON(w, q)
	if err != nil {
		log.Fatal(err)
	}
}

// All takes care of the 'all' route.
func All(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	err := writeJSON(w, repo.All())
	if err != nil {
		log.Fatal(err)
	}
}

// Search takes care of the /search/{word} route.
func Search(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	word := strings.ToLower(p.ByName("word"))
	qs := repo.AllByWord(word)

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
	q1 := repo.Random()
	q2 := repo.Random()

	quoteArray := strings.Split(q1.Quote, " ")
	quoteArray1 := strings.Split(q2.Quote, " ")
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

	q := repo.FindByUUID(uuidToSearch)
	if q == nil {
		err := writeNotFound(w, uuidToSearch)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	err := writeJSON(w, q)
	if err != nil {
		log.Fatal(err)
	}
}

// Like takes care of the /one/{UUID}/like route.
func Like(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	uuidToSearch := p.ByName("uuid")

	if err := repo.IncrementsScore(uuidToSearch); err != nil {
		err = writeNotFound(w, uuidToSearch)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
}

// Dislike takes care of the /one/{UUID}/dislike route.
func Dislike(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	uuidToSearch := p.ByName("uuid")

	if err := repo.DecrementsScore(uuidToSearch); err != nil {
		err = writeNotFound(w, uuidToSearch)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
}

// Top takes care of the /top route.
func Top(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	err := writeJSON(w, repo.Prefered())
	if err != nil {
		log.Fatal(err)
	}
}
