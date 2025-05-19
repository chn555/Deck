/*
The deck executable is a gRPC server that provides
a DeckService for managing decks of cards.

It listens on port 8080 and serves requests
for managing decks of cards.
*/
package main

import (
	"log/slog"
	"net"
	"os"

	deckImpl "github.com/chn555/deck/internal/deck/v1"
	"github.com/chn555/deck/pkg/deck"
	deckPb "github.com/chn555/schemas/proto/deck/v1"
	"google.golang.org/grpc"
)

func main() {
	// Create a new in-memory deck store so that decks can persist between requests
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
