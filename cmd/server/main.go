package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/yramanovich/word-of-wisdom/bucket"
	"github.com/yramanovich/word-of-wisdom/pow"
	"github.com/yramanovich/word-of-wisdom/quotes"
	"github.com/yramanovich/word-of-wisdom/tcp"
)

// Configure application with flags for the simplicity.
// In real application it's better to configure using environment variables.
var (
	profile             string
	port                string
	logLevel            string
	challengeExpiration time.Duration
)

func init() {
	flag.StringVar(&profile, "cpuprofile", "", "Write cpu profile to file")
	flag.StringVar(&port, "port", ":8080", "Server port")
	flag.StringVar(&logLevel, "loglevel", "info", "Logging level")
	flag.DurationVar(&challengeExpiration, "challenge-expiration", 3*time.Second,
		"Duration for which client should resolve the challenge")
	flag.Parse()
}

func run() int {
	logger := initLogger()

	if profile != "" {
		f, err := os.Create("cpu.prof")
		if err != nil {
			logger.Error("create cpu profile file", "err", err)
		} else {
			pprof.StartCPUProfile(f) // nolint:errcheck
			defer pprof.StopCPUProfile()
			defer f.Close() // nolint:errcheck
		}
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// Use embedded quote for the sake of simplicity.
	// If we want to add quotes without restarting the application,
	// we can create some db-oriented implementation with some API.
	embeddedQuoter := quotes.NewEmbeddedQuoter()

	// Use token bucket in order to increase the difficulty of challenger
	// if server faces some high load.
	tokenBucket := bucket.NewTokenBucket(1, 1000)

	// Filler difficulter increases difficulty of the challenge for every 100 points of overflow.
	// For n requests - difficulty 20, but for n + 100 difficulty 21.
	// Todo: it will be better to increase the complexity for each identified user separately.
	difficulter := pow.NewBucketDifficulter(tokenBucket, 100)

	// Challenges connection with hashcash proof-of-work algorithm.
	// Difficulter is used to increase difficulty dynamically according to the load.
	connectionChallenger := pow.NewHashcashChallenger(difficulter, logger, challengeExpiration)
	srv := tcp.NewServer(tcp.ServerConfig{
		Quoter:     embeddedQuoter,
		Port:       port,
		Logger:     logger,
		Challenger: connectionChallenger,
		Filler:     tokenBucket,
	})
	logger.Info("Start listening on the TCP socket", "port", port)
	defer logger.Info("Stop listening on the TCP socket", "port", port)

	if err := srv.Run(ctx); err != nil {
		logger.Error("run server", "err", err)
		return 1
	}

	return 0
}

func main() {
	os.Exit(run())
}

func initLogger() *slog.Logger {
	var lvl slog.Level
	_ = lvl.UnmarshalText([]byte(logLevel)) // nolint:errcheck
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))
}
