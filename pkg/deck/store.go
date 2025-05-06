package deck

import (
	"context"
	"fmt"
)

// Store is an abstraction of deck storage with id
type Store interface {
	// Get fetches the deck for the given id
	Get(ctx context.Context, id string) (*Deck, error)
	// Put stores the deck for the given id
	Put(ctx context.Context, id string, deck *Deck) error
}

// InMemoryDeckStore is a non-persistent implementation of Store
type InMemoryDeckStore struct {
	state map[string]*Deck
}

func NewInMemoryDeckStore() *InMemoryDeckStore {
	return &InMemoryDeckStore{
		state: make(map[string]*Deck),
	}
}

// Get fetches the deck for the given id
func (s *InMemoryDeckStore) Get(_ context.Context, id string) (*Deck, error) {
	deck, ok := s.state[id]
	if !ok {
		return nil, fmt.Errorf("deck %s not found", id)
	}

	return deck, nil
}

// Put stores the deck for the given id
func (s *InMemoryDeckStore) Put(_ context.Context, id string, deck *Deck) error {
	s.state[id] = deck
	return nil
}
