package handler

import (
	"math/rand"
	"strings"
	"time"
)

func stringModifier(minQuote []string, maxQuote []string) string {
	// Generate the number of words to be modified
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	swappable := r1.Intn(len(minQuote)-1) + 1

	for i := 0; i < swappable; i++ {
		si := rand.NewSource(time.Now().UnixNano() * 10000)
		ri := rand.New(si)

		// Choose a random place in each array and decide whether replacing or inserting
		minPlace := ri.Intn(len(minQuote)) * 7 % len(minQuote)
		maxPlace := ri.Intn(len(maxQuote)) * 9 % len(maxQuote)
		modifier := ri.Float64()

		// Replace or Insert
		if modifier > .5 {
			maxQuote[maxPlace] = minQuote[minPlace]
		} else {
			maxQuote = append(maxQuote, "")
			copy(maxQuote[maxPlace+1:], maxQuote[maxPlace:])
			maxQuote[maxPlace] = minQuote[minPlace]
		}
	}

	return strings.Join(maxQuote, " ")
}
