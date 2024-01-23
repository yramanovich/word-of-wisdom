package main

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/yramanovich/word-of-wisdom/bytesio"
	"github.com/yramanovich/word-of-wisdom/pow"
)

var (
	address  string
	logLevel string
)

func init() {
	flag.StringVar(&logLevel, "loglevel", "info", "Logging level")
	flag.StringVar(&address, "addr", "localhost:8080", "Server address")
	flag.Parse()
}

func run() int {
	slog.SetDefault(initLogger())

	conn, err := net.Dial("tcp", address)
	if err != nil {
		slog.Error("dial server", "err", err)
		return 1
	}
	defer conn.Close() // nolint:errcheck

	slog.Debug("waiting for challenge")

	challenge, err := bytesio.Readln(conn)
	if err != nil {
		slog.Error("read challenge", "err", err)
		return 1
	}

	slog.Debug("solve challenge from the server")

	start := time.Now()
	solved, err := pow.Solve(challenge)
	if err != nil {
		slog.Error("solve challenge", "err", err, "challenge", string(challenge))
		return 1
	}
	slog.Debug("solved challenge, send solution to server", "elapsed_time", time.Since(start))

	if err := bytesio.Writeln(conn, solved); err != nil {
		slog.Error("write solution to server", "err", err)
		return 1
	}

	quote, err := io.ReadAll(conn)
	if err != nil {
		slog.Error("receive quote from server", "err", err)
		return 1
	}

	fmt.Println(string(quote))
	return 0
}

func initLogger() *slog.Logger {
	var lvl slog.Level
	_ = lvl.UnmarshalText([]byte(logLevel)) // nolint:errcheck
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))
}

func main() {
	os.Exit(run())
}
