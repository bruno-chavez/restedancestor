package quotes

import (
	"testing"
)

type DataDouble struct{}

var raw = `[
  {
    "quote": "There is a place, beneath those ancient ruins, in the moor, that calls out to the boldest among them...",
    "uuid": "240f6e4e-3f28-4876-9ece-c88599c15c78",
    "score": 0
  },
  {
    "quote": "'We are the Flame!', they cry, 'And Darkness fears us!'",
    "uuid": "3d3bd030-74bf-4fcb-a0e0-b4dac38a688e",
    "score": 0
  }]`

// read fetches data from storage
func (f DataDouble) Read() []byte {
	return []byte(raw)
}

func (d DataDouble) Write(data []byte) error {
	raw = string(data[:])
	return nil // errors.New("unknown")
}

var db = &DataDouble{}
var quotes = Parser(*db)

func TestOffsetQuoteFromUUIDFound(t *testing.T) {
	uuidToSearch := "240f6e4e-3f28-4876-9ece-c88599c15c78"
	_, err := quotes.OffsetQuoteFromUUID(db, uuidToSearch)
	if err != nil {
		t.Errorf("Unknown offset")
	}
}

func TestOffsetQuoteFromUUIDNotFound(t *testing.T) {
	uuidToSearch := "not found"

	_, err := quotes.OffsetQuoteFromUUID(db, uuidToSearch)

	if err == nil {
		t.Errorf("Offset found : %s", uuidToSearch)
	}
}

func TestLikeQuote(t *testing.T) {
	uuidToSearch := "240f6e4e-3f28-4876-9ece-c88599c15c78"

	offset, err := quotes.OffsetQuoteFromUUID(db, uuidToSearch)
	if err != nil {
		t.Errorf("Error While fetching quote")
	}
	originalQuote := quotes[*offset]
	quotes.LikeQuote(db, uuidToSearch)
	newQuote := quotes[*offset]

	if originalQuote.Score+1 != newQuote.Score {
		t.Errorf("New score isn't equal to original score + 1")
	}
}

func TestDislikeQuote(t *testing.T) {
	uuidToSearch := "3d3bd030-74bf-4fcb-a0e0-b4dac38a688e"

	offset, err := quotes.OffsetQuoteFromUUID(db, uuidToSearch)
	if err != nil {
		t.Errorf("Error While fetching quote")
	}
	originalQuote := quotes[*offset]
	quotes.DislikeQuote(db, uuidToSearch)
	newQuote := quotes[*offset]

	if originalQuote.Score-1 != newQuote.Score {
		t.Errorf("New score isn't equal to original score - 1")
	}
}
