package bucket

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testTimer struct {
	since time.Duration
}

func (*testTimer) Now() time.Time                  { return time.Now() }
func (t *testTimer) Since(time.Time) time.Duration { return t.since }

func TestTokenBucket_Fill(t *testing.T) {
	tb := NewTokenBucket(1, 2)
	ttimer := &testTimer{}
	tb.timer = ttimer

	// Add 2 tokens to the bucket -> no tokens left
	assert.Equal(t, 0, tb.Fill(2))

	// Add 1 more token to the bucket -> bucket going beyond for 1 token
	assert.Equal(t, 1, tb.Fill(1))

	// Wait 2 seconds with rate 1 -> add 2 tokens to the bucket
	ttimer.since = 2 * time.Second

	// One token should be available at this point
	assert.Equal(t, 0, tb.Fill(1))

}
