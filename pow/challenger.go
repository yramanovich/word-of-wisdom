package pow

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/yramanovich/world-of-wisdom/bytesio"
)

// Difficulter provides challenge difficulty for HashcashChallenger.
type Difficulter interface {
	Difficulty() int
}

// HashcashChallenger is an implementation for hash-based proof-of-work algorithm,
// which challenges net.Conn.
type HashcashChallenger struct {
	logger      *slog.Logger
	difficulter Difficulter
	expDuration time.Duration
}

// NewHashcashChallenger return new HashcashChallenger instance.
func NewHashcashChallenger(diff Difficulter, logger *slog.Logger, exp time.Duration) *HashcashChallenger {
	return &HashcashChallenger{
		logger:      logger,
		difficulter: diff,
		expDuration: exp,
	}
}

// Challenge challenges connection with hashcash algorithm.
func (c HashcashChallenger) Challenge(_ context.Context, conn net.Conn) error {
	difficulty := c.difficulter.Difficulty()
	resource := []byte(conn.RemoteAddr().String())

	c.logger.Debug("generate new hashcash", "rsc", string(resource), "difficulty", difficulty)

	start := time.Now()
	hash := New(resource, difficulty, start)

	c.logger.Debug("challenge client with hash")

	if err := bytesio.Writeln(conn, hash); err != nil {
		return fmt.Errorf("write hash to connection: %w", err)
	}

	c.logger.Debug("waiting for solution from client")

	solution, err := bytesio.Readln(conn)
	if err != nil {
		return fmt.Errorf("read solution: %w", err)
	}

	c.logger.Debug("verify solution from the client")

	if err := Verify(solution, hash, c.expDuration); err != nil {
		return fmt.Errorf("verify hashcash: %w", err)
	}

	c.logger.Debug("given solution passes all tests, it is considered a valid hash string")

	return nil
}
