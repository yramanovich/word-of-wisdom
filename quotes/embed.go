package quotes

import (
	"context"
	_ "embed"
	"math/rand"
	"strings"
)

//go:embed quotes.txt
var quotesContent string

// EmbeddedQuoter holds list of quotes in memory.
type EmbeddedQuoter struct {
	quotes   []string
	randomer func(int) int
}

// NewEmbeddedQuoter returns new EmbeddedQuoter instance.
func NewEmbeddedQuoter() *EmbeddedQuoter {
	split := strings.Split(quotesContent, "\n")
	quotes := make([]string, 0, len(split))
	for _, s := range split {
		s = strings.TrimSpace(s)
		if s != "" {
			quotes = append(quotes, s)
		}
	}

	return &EmbeddedQuoter{
		quotes:   quotes,
		randomer: rand.Intn,
	}
}

// GetQuote returns one random quote from the list.
func (es *EmbeddedQuoter) GetQuote(_ context.Context) string {
	num := es.randomer(len(es.quotes))
	return es.quotes[num]
}
