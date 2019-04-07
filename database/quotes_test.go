package database

import (
	"testing"
)

// à but de test, on peut faire un structure file qui a pour fonction : read et write
// On injecte cette structure dans le parser pour qu'elle soit lue et travailler en déféré
// db already exists in database.go
var quotes = Parser(*db)

func TestOffsetQuoteFromUUIDFound(t *testing.T) {
	uuidToSearch := "240f6e4e-3f28-4876-9ece-c88599c15c78"
	_, err := quotes.OffsetQuoteFromUUID(uuidToSearch)
	if err != nil {
		t.Errorf("Unknown offset")
	}
}

func TestOffsetQuoteFromUUIDNotFound(t *testing.T) {
	uuidToSearch := "240f6e4e-3f28-4876-9ece-c885995c78"

	_, err := quotes.OffsetQuoteFromUUID(uuidToSearch)

	if err == nil {
		t.Errorf("Offset found : %s", uuidToSearch)
	}
}

func TestLikeQuote(t *testing.T) {
	uuidToSearch := "240f6e4e-3f28-4876-9ece-c88599c15c78"

	offset, err := quotes.OffsetQuoteFromUUID(uuidToSearch)
	if err != nil {
		t.Errorf("Error While fetching quote")
	}
	originalQuote := quotes[*offset]
	quotes.LikeQuote(uuidToSearch)
	newQuote := quotes[*offset]

	if originalQuote.Score+1 != newQuote.Score {
		t.Errorf("New score isn't equal to original score + 1")
	}
}

func TestDislikeQuote(t *testing.T) {
	uuidToSearch := "45bdaed9-6147-4897-bbd9-fe223785a918"

	offset, err := quotes.OffsetQuoteFromUUID(uuidToSearch)
	if err != nil {
		t.Errorf("Error While fetching quote")
	}
	originalQuote := quotes[*offset]
	quotes.DislikeQuote(uuidToSearch)
	newQuote := quotes[*offset]

	if originalQuote.Score-1 != newQuote.Score {
		t.Errorf("New score isn't equal to original score - 1")
	}
}
