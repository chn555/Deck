package v1

import (
	"context"
	"fmt"

	deckPb "github.com/chn555/deck/internal/gen/proto/deck/v1"
	"github.com/chn555/deck/pkg/deck"
	"github.com/google/uuid"
)

type DeckServiceServer struct {
	store deck.Store
	deckPb.UnimplementedDeckServiceServer
}

func NewDeckServiceServer(store deck.Store) (*DeckServiceServer, error) {
	if store == nil {
		return nil, fmt.Errorf("store is nil")
	}
	return &DeckServiceServer{store: store}, nil
}

func (s DeckServiceServer) FetchCard(ctx context.Context, request *deckPb.FetchCardRequest) (*deckPb.Card, error) {
	d, err := s.store.Get(ctx, request.GetDeck().GetDeckId())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch deck: %w", err)
	}

	card, ok := d.FetchCard()
	if !ok {
		return nil, fmt.Errorf("failed to fetch card, deck is empty")
	}

	return cardToProto(card), nil
}

func (s DeckServiceServer) PushCard(ctx context.Context, request *deckPb.PushCardRequest) (*deckPb.Empty, error) {
	d, err := s.store.Get(ctx, request.GetDeck().GetDeckId())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch deck: %w", err)
	}

	d.PushCard(protoToCard(request.GetCard()))
	return &deckPb.Empty{}, nil
}

func (s DeckServiceServer) Create(ctx context.Context, _ *deckPb.Empty) (*deckPb.Deck, error) {
	newDeck := deck.NewDeck(
		deck.WithJokersInDeck(2),
		deck.WithShuffle(true),
	)
	id := uuid.New().String()
	err := s.store.Put(ctx, id, newDeck)
	if err != nil {
		return nil, err
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
	case deck.Aces:
		return deckPb.Card_SUIT_ACES
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
	case deckPb.Card_SUIT_ACES:
		return deck.Aces
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
