package quotes

import (
	"testing"
)

type DataDouble struct{}

var raw = `{
  "data": [
    {
      "quote": "beneath those ancient ruins,...",
      "uuid": "240f6e4e-3f28-4876-9ece-c88599c15c78",
      "score": 0
    },
    {
      "quote": "'We are the Ruins!'",
      "uuid": "3d3bd030-74bf-4fcb-a0e0-b4dac38a688e",
      "score": 0
    }
  ],
  "index": [
    {
      "word": "beneath",
      "uuids": [
        "240f6e4e-3f28-4876-9ece-c88599c15c78"
      ]
    },
    {
      "word": "those",
      "uuids": [
        "240f6e4e-3f28-4876-9ece-c88599c15c78"
      ]
    },
    {
      "word": "ancient",
      "uuids": [
        "240f6e4e-3f28-4876-9ece-c88599c15c78"
      ]
    },
    {
      "word": "ruins",
      "uuids": [
        "240f6e4e-3f28-4876-9ece-c88599c15c78",
        "3d3bd030-74bf-4fcb-a0e0-b4dac38a688e"
      ]
    }
  ]
}`

// Read provides in-memory data
func (d DataDouble) Read() []byte {
	return []byte(raw)
}

// Write stores new data into in-memory test double
func (d DataDouble) Write(data []byte) error {
	raw = string(data[:])
	return nil // errors.New("unknown")
}

var db = &DataDouble{}
var quotes = Parser(*db)

// TestLen tests only one part of the Len function
// TODO : test other part of the Len definition domain, ie. incrementation and decrementation
func TestLen(t *testing.T) {
	if quotes.Len() != 2 {
		t.Error("Length is not equal")
	}
}
func TestSwap(t *testing.T) {
	old1Pos, _ := quotes.OffsetQuoteFromUUID("240f6e4e-3f28-4876-9ece-c88599c15c78")
	old2Pos, _ := quotes.OffsetQuoteFromUUID("3d3bd030-74bf-4fcb-a0e0-b4dac38a688e")
	quotes.Swap(*old1Pos, *old2Pos)

	new1Pos, _ := quotes.OffsetQuoteFromUUID("240f6e4e-3f28-4876-9ece-c88599c15c78")
	new2Pos, _ := quotes.OffsetQuoteFromUUID("3d3bd030-74bf-4fcb-a0e0-b4dac38a688e")

	if *old1Pos != *new2Pos || *old2Pos != *new1Pos {
		t.Error("Position didn't swap")
	}
}

func TestLess(t *testing.T) {
	uuid := "240f6e4e-3f28-4876-9ece-c88599c15c78"
	element1Pos, _ := quotes.OffsetQuoteFromUUID(uuid)
	element2Pos, _ := quotes.OffsetQuoteFromUUID("3d3bd030-74bf-4fcb-a0e0-b4dac38a688e")

	if quotes.Less(*element1Pos, *element2Pos) {
		t.Error("Elements are sorted")
	}

	quotes.LikeQuote(db, uuid)

	if !quotes.Less(*element1Pos, *element2Pos) {
		t.Error("Elements are equals")
	}

}

// TestOffsetQuoteFromUUIDFound checks a UUID can be found
func TestOffsetQuoteFromUUIDFound(t *testing.T) {
	uuidToSearch := "240f6e4e-3f28-4876-9ece-c88599c15c78"
	_, err := quotes.OffsetQuoteFromUUID(uuidToSearch)
	if err != nil {
		t.Errorf("Unknown offset")
	}
}

// TestOffsetQuoteFromUUIDNotFound checks non existent UUID returns error
func TestOffsetQuoteFromUUIDNotFound(t *testing.T) {
	uuidToSearch := "not found"

	_, err := quotes.OffsetQuoteFromUUID(uuidToSearch)

	if err == nil {
		t.Errorf("Offset found : %s", uuidToSearch)
	}
}

// TestLikeQuote checks the quote's score increments
func TestLikeQuote(t *testing.T) {
	uuidToSearch := "240f6e4e-3f28-4876-9ece-c88599c15c78"

	offset, err := quotes.OffsetQuoteFromUUID(uuidToSearch)
	if err != nil {
		t.Errorf("Error While fetching quote")
	}
	originalQuote := quotes.Data[*offset]
	quotes.LikeQuote(db, uuidToSearch)
	newQuote := quotes.Data[*offset]

	if originalQuote.Score+1 != newQuote.Score {
		t.Errorf("New score isn't equal to original score + 1")
	}
}

// TestDislikeQuote checks the quote's score decrements
func TestDislikeQuote(t *testing.T) {
	uuidToSearch := "3d3bd030-74bf-4fcb-a0e0-b4dac38a688e"

	offset, err := quotes.OffsetQuoteFromUUID(uuidToSearch)
	if err != nil {
		t.Errorf("Error While fetching quote")
	}
	originalQuote := quotes.Data[*offset]
	quotes.DislikeQuote(db, uuidToSearch)
	newQuote := quotes.Data[*offset]

	if originalQuote.Score-1 != newQuote.Score {
		t.Errorf("New score isn't equal to original score - 1")
	}
}

func TestIndex(t *testing.T) {
	originalQuotes := quotes
	originalIndexes := originalQuotes.Indexes

	workingQuotes := quotes
	workingQuotes.Index(db)
	workingIndexes := workingQuotes.Indexes

	if len(workingIndexes) != len(originalIndexes) {
		t.Errorf("Quotes Lengths aren't equals, %d <> %d", len(workingIndexes), len(originalIndexes))
	}

	for i, workingIndex := range originalIndexes {
		originalIndex := originalIndexes[i]
		if len(workingIndex.Uuids) != len(originalIndex.Uuids) {
			t.Errorf("For # %d, %s and %s are different", i, originalIndexes[i], workingIndexes[i])
		}
	}
}

func TestListFilled(t *testing.T) {
	list := quotes.List("ruins")

	if len(list) != 2 {
		t.Errorf("ruins wasn't found in %d documents", len(list))
	}
}

func TestListEmpty(t *testing.T) {
	list := quotes.List("dfzefzer")

	if len(list) != 0 {
		t.Errorf("random text was found in documents")
	}
}
