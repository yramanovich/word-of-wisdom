package pow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummyBucket int

func (d dummyBucket) Fill(int) int { return int(d) }

func TestBucketDifficultyResolver(t *testing.T) {
	assert.Equal(t, defaultDifficulty, NewBucketDifficulter(dummyBucket(10), 100).Difficulty())
	assert.Equal(t, 22, NewBucketDifficulter(dummyBucket(250), 100).Difficulty())
}
