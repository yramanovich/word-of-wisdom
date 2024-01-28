package pow

import (
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHashcash(t *testing.T) {
	start := time.Now()
	hc := New([]byte("new"), 10, start)
	assert.Regexp(t, regexp.MustCompile("^1:10:\\d*:\\S+:\\S+$"), string(hc))

	solved, err := Solve(hc)
	assert.NoError(t, err)
	assert.Regexp(t, regexp.MustCompile("^1:10:\\d*:\\S+:\\S+:\\S+$"), string(solved))

	// for better testing we should provide time as interface, but im lazy
	err = Verify(solved, hc, 20*time.Second)
	assert.NoError(t, err)

	// test expiration
	time.Sleep(20 * time.Millisecond)
	err = Verify(solved, hc, time.Nanosecond)
	assert.Error(t, err)

	data, err := Solve([]byte("invalidhash"))
	assert.Error(t, err)
	assert.Empty(t, data)

	data, err = Solve(nil)
	assert.Error(t, err)
	assert.Empty(t, data)

	err = Verify([]byte("1:0:1705962838:bmV3:Bod8ozrlUcqUImyR2swBuQ==:SwAAAA=="), hc, time.Second)
	assert.Error(t, err)
}

func Test_hasNLeadingZeroes(t *testing.T) {
	type args struct {
		hash []byte
		n    int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "has 10 leading zeroes",
			args: args{
				hash: []byte{0x0, 0x0},
				n:    10,
			},
			want: true,
		},
		{
			name: "has not 10 leading zeroes",
			args: args{
				hash: []byte{0x0, 0x50},
				n:    10,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, hasNLeadingZeroes(tt.args.hash, tt.args.n),
				"hasNLeadingZeroes(%v, %v)", tt.args.hash, tt.args.n)
		})
	}
}
