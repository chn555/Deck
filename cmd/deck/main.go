package main

import (
	"context"
	"log/slog"
	"net"
	"os"

	deckImpl "github.com/chn555/deck/internal/deck/v1"
	"github.com/chn555/deck/pkg/deck"
	deckPb "github.com/chn555/schemas/proto/deck/v1"
	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func main() {
	store := deck.NewInMemoryDeckStore()
	deckServer, err := deckImpl.NewDeckServiceServer(store)
	if err != nil {
		slog.Error("Error creating deck service", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	opts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
		// Add any other option (check functions starting with logging.With).
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(InterceptorLogger(logger), opts...),
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(InterceptorLogger(logger), opts...),
		),
	)
	deckPb.RegisterDeckServiceServer(grpcServer, deckServer)
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		slog.Error("Error creating listener", err)
		os.Exit(1)
	}

	slog.Info("Listening on port 8080")
	err = grpcServer.Serve(listener)
	if err != nil {
		slog.Error("Error serving grpc server", err)
		os.Exit(1)
	}
}

// InterceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
