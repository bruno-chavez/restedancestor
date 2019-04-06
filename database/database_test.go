package database

import (
	"testing"
)

var database = Parser()

func TestOffsetQuoteFromUUIDFound(t *testing.T) {
	uuidToSearch := "240f6e4e-3f28-4876-9ece-c88599c15c78"
	_, err := OffsetQuoteFromUUID(uuidToSearch)
	if err != nil {
		t.Errorf("Unknown offset")
	}
}

func TestOffsetQuoteFromUUIDNotFound(t *testing.T) {
	uuidToSearch := "240f6e4e-3f28-4876-9ece-c885995c78"

	_, err := OffsetQuoteFromUUID(uuidToSearch)
	if err == nil {
		t.Errorf("Offset found : %s", uuidToSearch)
	}
}

func TestLikeQuote(t *testing.T) {
	uuidToSearch := "240f6e4e-3f28-4876-9ece-c88599c15c78"

	offset, err := OffsetQuoteFromUUID(uuidToSearch)
	if err != nil {
		t.Errorf("Error While fetching quote")
	}
	originalQuote := database[*offset]
	LikeQuote(uuidToSearch)
	newQuote := Parser()[*offset]

	if originalQuote.Score+1 != newQuote.Score {
		t.Errorf("New score isn't equal to original score + 1")
	}
}

func TestDislikeQuote(t *testing.T) {
	uuidToSearch := "45bdaed9-6147-4897-bbd9-fe223785a918"

	offset, err := OffsetQuoteFromUUID(uuidToSearch)
	if err != nil {
		t.Errorf("Error While fetching quote")
	}
	originalQuote := database[*offset]
	DislikeQuote(uuidToSearch)
	newQuote := Parser()[*offset]

	if originalQuote.Score-1 != newQuote.Score {
		t.Errorf("New score isn't equal to original score - 1")
	}
}
