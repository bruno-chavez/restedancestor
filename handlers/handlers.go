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

// Random handles the '/random' endpoint
func Random(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {

	q := repo.Random()
	err := writeResponse(w, q)
	if err != nil {
		log.Fatal(err)
	}
}

// All handles the '/all' endpoint
func All(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {

	err := writeResponse(w, repo.All())
	if err != nil {
		log.Fatal(err)
	}
}

// Search handles the '/search' endpoint
func Search(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {

	word := strings.ToLower(p.ByName("word"))
	qs := repo.AllByWord(word)

	if len(qs) != 0 {
		err := writeResponse(w, qs)
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

// Senile handles the '/senile' endpoint
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

	joinedQuote := quotes.Quote{Quote: quote}
	err := writeResponse(w, joinedQuote)
	if err != nil {
		log.Fatal(err)
	}
}

// Find handles the '/uuid/:uuid/find' endpoint
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

	err := writeResponse(w, q)
	if err != nil {
		log.Fatal(err)
	}
}

// Like handles the '/uuid/:uuid/like' endpoint
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

// Dislike handles the '/uuid/:uuid/dislike' endpoint
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

// Top handles the '/top' endpoint
func Top(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {

	err := writeResponse(w, repo.Preferred())
	if err != nil {
		log.Fatal(err)
	}
}
