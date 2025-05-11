package deck

import (
	"context"
	"fmt"
)

// Store is an abstraction of deck storage with id.
type Store interface {
	// Get fetches the deck for the given id.
	Get(ctx context.Context, id string) (*Deck, error)
	// Put stores the deck for the given id.
	Put(ctx context.Context, id string, deck *Deck) error
}

// InMemoryDeckStore is an implementation of Store that stores decks in memory.
//
// It is provided here as a simplification of an actual database implementation
// to not complicate this exercise with the actual concerns of
// working with a database.
//
// Since no data is being written to disk, this store will not persist after a
// restart of the application.
type InMemoryDeckStore struct {
	state map[string]*Deck
}

// NewInMemoryDeckStore creates a new InMemoryDeckStore.
func NewInMemoryDeckStore() *InMemoryDeckStore {
	return &InMemoryDeckStore{
		state: make(map[string]*Deck),
	}
}

// Get fetches the deck for the given id.
func (s *InMemoryDeckStore) Get(_ context.Context, id string) (*Deck, error) {
	deck, ok := s.state[id]
	if !ok {
		return nil, fmt.Errorf("deck not found")
	}

	return deck, nil
}

// Put stores the deck for the given id.
func (s *InMemoryDeckStore) Put(_ context.Context, id string, deck *Deck) error {
	s.state[id] = deck
	return nil
}
