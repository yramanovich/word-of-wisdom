package quotes

import (
	"context"
	"strings"
	"testing"
)

func TestNewEmbeddedQuoter(t *testing.T) {
	{
		// test with deterministic randomer
		q := NewEmbeddedQuoter()
		q.randomer = func(int) int { return 0 }

		got := q.GetQuote(context.TODO())
		if got != `"The only thing we have to fear is fear itself." - Franklin D. Roosevelt` {
			t.Error("quotes do not match")
		}
	}

	{
		// test with non-deterministic randomer
		q := NewEmbeddedQuoter()
		got := q.GetQuote(context.TODO())
		if !strings.Contains(quotesContent, got) {
			t.Error("quotes content missing got")
		}
	}
}
