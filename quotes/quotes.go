package quotes

import (
	"github.com/satori/go.uuid"
)

// QuoteType describes a quote.
type QuoteType struct {
	ID    int       `json:"id"`
	Quote string    `json:"quote"`
	UUID  uuid.UUID `json:"uuid"`
	Score int       `json:"score"`
}
