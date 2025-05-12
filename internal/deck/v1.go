package deck

import (
	"context"

	deckPb "github.com/chn555/schemas/proto/deck/v1"
)

// V1DeckServiceServer is an implementation of the DeckServiceServer interface.
// It contains the implementation of the methods defined in the proto file.
type V1DeckServiceServer struct {
	// UnimplementedDeckServiceServer is embedded to ensure forward compatibility,
	// and to provide a default implementation of the methods.
	deckPb.UnimplementedDeckServiceServer
}

// FetchCard fetches a card for the given deck id.
func (s *V1DeckServiceServer) FetchCard(ctx context.Context, request *deckPb.FetchCardRequest) (*deckPb.Card, error) {
	return s.UnimplementedDeckServiceServer.FetchCard(ctx, request)
}

// PushCard pushes a card to the given deck id.
func (s *V1DeckServiceServer) PushCard(ctx context.Context, request *deckPb.PushCardRequest) (*deckPb.Empty, error) {
	return s.UnimplementedDeckServiceServer.PushCard(ctx, request)
}

// Create creates a new deck with the given configuration.
func (s *V1DeckServiceServer) Create(ctx context.Context, req *deckPb.CreateDeckRequest) (*deckPb.Deck, error) {
	return s.UnimplementedDeckServiceServer.Create(ctx, req)
}
