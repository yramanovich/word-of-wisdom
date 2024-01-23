package pow

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strconv"
	"time"
)

const (
	version = 1
	delim   = ':'
)

type hashcash struct {
	version  int    // Hashcash format version.
	bits     int    // Number of "partial pre-image" (zero) bits in the hashed code.
	date     int64  // The time that the message was sent, in unix format.
	resource []byte // Resource data (an IP address).
	nonce    []byte // Random characters, encoded in base-64 format.
	counter  []byte // Binary counter, encoded in base-64 format.
}

func (hc hashcash) encoded() []byte {
	ver := strconv.Itoa(hc.version)
	bits := strconv.Itoa(hc.bits)
	date := strconv.FormatInt(hc.date, 10)

	toJoin := [][]byte{
		[]byte(ver),
		[]byte(bits),
		[]byte(date),
		hc.resource,
		hc.nonce,
	}

	if len(hc.counter) != 0 {
		toJoin = append(toJoin, hc.counter)
	}

	return bytes.Join(toJoin, []byte{delim})
}

// New returns new hashcash stamp.
func New(rsc []byte, difficulty int, date time.Time) []byte {
	if len(rsc) == 0 {
		panic("rsc can't be empty")
	}

	if difficulty <= 0 {
		panic("difficulty has to be > 0")
	}

	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		panic(err) // Omit error handling for simplicity
	}

	encodedNonce := make([]byte, base64.URLEncoding.EncodedLen(len(nonce)))
	base64.URLEncoding.Encode(encodedNonce, nonce)

	encodedResource := make([]byte, base64.URLEncoding.EncodedLen(len(rsc)))
	base64.URLEncoding.Encode(encodedResource, rsc)

	hc := hashcash{
		version:  version,
		bits:     difficulty,
		date:     date.Unix(),
		resource: encodedResource,
		nonce:    encodedNonce,
	}
	return hc.encoded()
}

// Solve finds solution for the given hashcash stamp.
func Solve(b []byte) ([]byte, error) {
	hc, err := parse(b)
	if err != nil {
		return nil, fmt.Errorf("parse hashcash: %w", err)
	}

	n := countBits(hc.bits)

	solution := uint32(0)
	sb := make([]byte, 4)
	for {
		binary.LittleEndian.PutUint32(sb, solution)

		encodedSolution := make([]byte, base64.URLEncoding.EncodedLen(len(sb)))
		base64.URLEncoding.Encode(encodedSolution, sb)

		hc.counter = encodedSolution

		hash := sha256.Sum256(bytes.Join([][]byte{b, encodedSolution}, []byte{delim}))
		if !hasNBits(hash[:n], hc.bits, n) {
			solution++
			continue
		}
		break
	}

	return hc.encoded(), nil
}

// Verify verifies that the stamp is valid.
func Verify(solution, challenge []byte, expiration time.Duration) error {
	if !bytes.HasPrefix(solution, challenge) {
		return fmt.Errorf("invalid solution: %s", string(solution))
	}
	hc, err := parse(solution)
	if err != nil {
		return err
	}

	date := time.Unix(hc.date, 0)
	if time.Since(date) > expiration {
		return fmt.Errorf("expired hashcash: %s", time.Since(date).String())
	}

	hash := sha256.Sum256(solution)
	n := countBits(hc.bits)

	if !hasNBits(hash[:n], hc.bits, n) {
		return fmt.Errorf("invalid solution")
	}

	return nil
}

func hasNBits(hash []byte, bits, n int) bool {
	if bits == 0 {
		return true
	}

	for i := 0; i < n; i++ {
		if bits > 8 {
			bits -= 8
			if hash[i] != 0 {
				return false
			}
			continue
		}

		// (bits % 8) == bits
		pad := 8 - bits
		if hash[i]>>pad == 0 {
			return true
		}
	}
	return false
}

func parse(b []byte) (hashcash, error) {
	var hc hashcash
	split := bytes.Split(b, []byte{delim})
	if !(len(split) == 5 || len(split) == 6) {
		return hc, fmt.Errorf("invalid count of segments")
	}

	ver, err := strconv.Atoi(string(split[0]))
	if err != nil {
		return hc, fmt.Errorf("invalid version: %w", err)
	}
	hc.version = ver

	bits, err := strconv.Atoi(string(split[1]))
	if err != nil {
		return hc, fmt.Errorf("invalid bits: %w", err)
	}
	hc.bits = bits

	date, err := strconv.ParseInt(string(split[2]), 10, 64)
	if err != nil {
		return hc, fmt.Errorf("invalid date: %w", err)
	}
	hc.date = date

	resource := split[3]
	if len(resource) == 0 {
		return hc, fmt.Errorf("invalid resource: %q", resource)
	}
	hc.resource = resource

	nonce := split[4]
	if len(nonce) == 0 {
		return hc, fmt.Errorf("invalid nonce: %q", nonce)
	}
	hc.nonce = nonce

	if len(split) > 5 {
		counter := split[5]
		if len(counter) == 0 {
			return hc, fmt.Errorf("invalid counter: %q", counter)
		}
		hc.counter = counter
	}

	return hc, nil
}

func countBits(bits int) int {
	n, m := bits/8, bits%8
	if m > 0 {
		n++
	}
	return n
}
