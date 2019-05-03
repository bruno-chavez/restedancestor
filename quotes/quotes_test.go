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
