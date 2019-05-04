package quotes

import (
	"github.com/satori/go.uuid"
)

// Quote describes a quote.
type Quote struct {
	id    int
	Quote string    `json:"quote"`
	UUID  uuid.UUID `json:"uuid"`
	Score int       `json:"score"`
}
