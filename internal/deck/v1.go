package deck

import (
	"context"
	"fmt"

	"github.com/chn555/deck/pkg/deck"
	deckPb "github.com/chn555/schemas/proto/deck/v1"
	"github.com/google/uuid"
)

// V1DeckServiceServer is an implementation of the DeckServiceServer interface.
// It contains the implementation of the methods defined in the proto file.
type V1DeckServiceServer struct {
	// store is the deck store
	store deck.Store

	// UnimplementedDeckServiceServer is embedded to ensure forward compatibility,
	// and to provide a default implementation of the methods.
	deckPb.UnimplementedDeckServiceServer
}

// NewV1DeckServiceServer creates a new V1DeckServiceServer.
func NewV1DeckServiceServer(store deck.Store) (*V1DeckServiceServer, error) {
	if store == nil {
		return nil, fmt.Errorf("store is nil")
	}
	return &V1DeckServiceServer{store: store}, nil
}

// FetchCard fetches a card for the given deck id. Once the card is fetched, it
// is removed from the deck and the modified deck is stored.
func (s *V1DeckServiceServer) FetchCard(ctx context.Context, request *deckPb.FetchCardRequest) (*deckPb.Card, error) {
	d, err := s.store.Get(ctx, request.GetDeck().GetDeckId())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch deck: %w", err)
	}

	card, ok := d.FetchCard()
	if !ok {
		return nil, fmt.Errorf("failed to fetch card, deck is empty")
	}

	err = s.store.Put(ctx, request.GetDeck().GetDeckId(), d)
	if err != nil {
		return nil, fmt.Errorf("failed to store deck: %w", err)
	}

	return cardToProto(card), nil
}

// PushCard pushes a card to the given deck id.
// Once the card is pushed, the modified deck is stored.
func (s *V1DeckServiceServer) PushCard(ctx context.Context, request *deckPb.PushCardRequest) (*deckPb.Empty, error) {
	d, err := s.store.Get(ctx, request.GetDeck().GetDeckId())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch deck: %w", err)
	}

	d.PushCard(protoToCard(request.GetCard()))

	err = s.store.Put(ctx, request.GetDeck().GetDeckId(), d)
	if err != nil {
		return nil, fmt.Errorf("failed to store deck: %w", err)
	}

	return &deckPb.Empty{}, nil
}

// Create creates a new deck with the given configuration and stores it.
func (s *V1DeckServiceServer) Create(ctx context.Context, req *deckPb.CreateDeckRequest) (*deckPb.Deck, error) {
	newDeck := deck.NewDeck(
		deck.WithJokersInDeck(int(req.GetJokerCount())),
		deck.WithShuffle(req.GetShuffle()),
		deck.WithAdditionalDecks(int(req.AdditionalDeckCount)),
	)
	id := uuid.New().String()
	err := s.store.Put(ctx, id, newDeck)
	if err != nil {
		return nil, fmt.Errorf("failed to store deck: %w", err)
	}

	return &deckPb.Deck{DeckId: id}, nil
}

func cardToProto(card deck.Card) *deckPb.Card {
	return &deckPb.Card{
		Value: int32(card.Value),
		Suit:  cardSuitToProto(card),
	}
}

func cardSuitToProto(card deck.Card) deckPb.Card_Suit {
	switch card.Suit {
	case deck.Hearts:
		return deckPb.Card_SUIT_HEARTS
	case deck.Spades:
		return deckPb.Card_SUIT_SPADES
	case deck.Diamonds:
		return deckPb.Card_SUIT_DIAMONDS
	case deck.Clubs:
		return deckPb.Card_SUIT_CLUBS
	case deck.Jokers:
		return deckPb.Card_SUIT_JOKERS
	default:
		return deckPb.Card_SUIT_UNSPECIFIED
	}
}

func protoToCard(card *deckPb.Card) deck.Card {
	return deck.NewCard(protoSuitToCard(card), uint8(card.Value))
}

func protoSuitToCard(card *deckPb.Card) deck.Suit {
	switch card.Suit {
	case deckPb.Card_SUIT_HEARTS:
		return deck.Hearts
	case deckPb.Card_SUIT_SPADES:
		return deck.Spades
	case deckPb.Card_SUIT_CLUBS:
		return deck.Clubs
	case deckPb.Card_SUIT_JOKERS:
		return deck.Jokers
	default:
		return deck.UnknownSuit
	}
}
