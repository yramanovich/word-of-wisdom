package tcp

import (
	"context"
	"log/slog"
	"net"
	"os"
	"sync"
	"time"
)

const (
	defaultPort = ":8080"
)

// Filler fills filler with tokens. Used to track server load.
type Filler interface {
	Fill(token int) int
}

// Challenger challenges connection.
type Challenger interface {
	Challenge(ctx context.Context, conn net.Conn) error
}

// Quoter returns rand quotes.
type Quoter interface {
	GetQuote(ctx context.Context) (quote string)
}

// ServerConfig is configuration for NewServer.
type ServerConfig struct {
	Quoter     Quoter
	Port       string
	Logger     *slog.Logger
	Challenger Challenger
	Filler     Filler
}

// NewServer returns new Server instance.
func NewServer(c ServerConfig) *Server {
	if c.Quoter == nil {
		panic("missing quoter parameter")
	}
	if c.Challenger == nil {
		panic("missing challenger parameter")
	}
	if c.Filler == nil {
		panic("missing filler parameter")
	}
	if c.Port == "" {
		c.Port = defaultPort
	}
	if c.Logger == nil {
		c.Logger = slog.Default()
	}

	return &Server{
		port:       c.Port,
		quoter:     c.Quoter,
		logger:     c.Logger,
		wg:         &sync.WaitGroup{},
		challenger: c.Challenger,
		filler:     c.Filler,
	}
}

// Server is a general component of the service.
// Challenges incoming connections and gives nice quotes if client does some work.
type Server struct {
	challenger Challenger
	filler     Filler
	port       string
	quoter     Quoter
	logger     *slog.Logger
	wg         *sync.WaitGroup
}

// Run starts tcp server.
func (s *Server) Run(ctx context.Context) error {
	addr, err := net.ResolveTCPAddr("tcp", s.port)
	if err != nil {
		return err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	defer l.Close() // nolint:errcheck

	defer s.wg.Wait()

	for {
		select {
		case <-ctx.Done():
			return nil

		default:
			if err := l.SetDeadline(time.Now().Add(200 * time.Millisecond)); err != nil {
				return err
			}

			conn, err := l.Accept()
			if err != nil {
				// log if not timeout error
				if !os.IsTimeout(err) {
					slog.Error("accept connection", "err", err)
				}
				continue
			}

			s.wg.Add(1)
			go s.handleConnection(ctx, conn)
		}
	}
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close() // nolint:errcheck

	// Add connection to the token filler.
	s.filler.Fill(1)

	if err := conn.SetDeadline(time.Now().Add(10 * time.Second)); err != nil {
		s.logger.Error("set connection deadline", "err", err)
		return
	}

	s.logger.Debug("handle connection", "address", conn.RemoteAddr().String())

	if err := s.challenger.Challenge(ctx, conn); err != nil {
		s.logger.Warn("client failed challenge", "err", err)
		return
	}

	quote := s.quoter.GetQuote(ctx)
	if _, err := conn.Write([]byte(quote)); err != nil {
		s.logger.Error("write quote", "err", err)
	}
}
