package pow

import (
	"bytes"
	"context"
	"log/slog"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHashcashChallenger_Challenge(t *testing.T) {
	challenger := NewHashcashChallenger(dummyDifficulter(1), slog.New(discardHandler{}), 10*time.Second)

	// todo: use tool like https://github.com/vektra/mockery to provide mocking
	conn := &dummyConn{}
	err := challenger.Challenge(context.TODO(), conn)
	assert.NoError(t, err)

}

type dummyDifficulter int

func (d dummyDifficulter) Difficulty() int { return int(d) }

type dummyConn struct {
	solution []byte
}

func (d *dummyConn) Read(b []byte) (n int, err error) {
	return bytes.NewReader(append(d.solution, '\n')).Read(b)
}

func (d *dummyConn) Write(b []byte) (n int, err error) {
	solution, _ := Solve(b[:len(b)-1])
	d.solution = solution
	return len(b), nil
}

func (d *dummyConn) Close() error                     { return nil }
func (d *dummyConn) LocalAddr() net.Addr              { return &net.TCPAddr{Port: 8000} }
func (d *dummyConn) RemoteAddr() net.Addr             { return &net.TCPAddr{Port: 8000} }
func (d *dummyConn) SetDeadline(time.Time) error      { return nil }
func (d *dummyConn) SetReadDeadline(time.Time) error  { return nil }
func (d *dummyConn) SetWriteDeadline(time.Time) error { return nil }

type discardHandler struct{}

func (discardHandler) Enabled(context.Context, slog.Level) bool  { return false }
func (discardHandler) Handle(context.Context, slog.Record) error { return nil }
func (d discardHandler) WithAttrs([]slog.Attr) slog.Handler      { return d }
func (d discardHandler) WithGroup(string) slog.Handler           { return d }
