package main

import (
	"log/slog"
	"net"
	"os"

	deckImpl "github.com/chn555/deck/internal/deck/v1"
	deckPb "github.com/chn555/deck/internal/gen/proto/deck/v1"
	"github.com/chn555/deck/pkg/deck"
	"google.golang.org/grpc"
)

func main() {
	store := deck.NewInMemoryDeckStore()
	deckServer, err := deckImpl.NewDeckServiceServer(store)
	if err != nil {
		slog.Error("Error creating deck service", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
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
