package pow

const defaultDifficulty = 20

// Filler is used to fill bucket capacity.
type Filler interface {
	Fill(tokens int) int
}

// BucketDifficulter increases difficulty for every k amount of overflow.
type BucketDifficulter struct {
	tb Filler
	k  int
}

// NewBucketDifficulter returns new BucketDifficulter instance.
func NewBucketDifficulter(tb Filler, k int) *BucketDifficulter {
	if k == 0 {
		panic("k should be > 0")
	}
	return &BucketDifficulter{
		tb: tb,
		k:  k,
	}
}

// Difficulty returns difficulty according to bucket contents.
func (t BucketDifficulter) Difficulty() int {
	// We don't want to fill the bucket, we just need to understand whether it is overfilled and how much
	overfill := t.tb.Fill(0)
	return defaultDifficulty + overfill/t.k
}
