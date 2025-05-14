package deck

import (
	"math/rand"
	"slices"
)

// Deck represents a deck of cards.
type Deck struct {
	// Cards represents the cards in the deck.
	Cards []Card
}

// Config represents the configuration for the deck.
type Config struct {
	// cmpFunc is the function used to compare cards for sorting the deck.
	cmpFunc func(a, b Card) int
	// shuffleDeck indicates whether to shuffle the deck.
	shuffleDeck bool
	// jokerInDeckCount is the number of jokers in the deck.
	jokerInDeckCount int
	// excludeFunc is the function used to exclude cards from the deck.
	excludeFunc func(Card) bool
	// additionalDecksCount is the number of additional decks to create.
	additionalDecksCount int
}

// ConfigOption is a function that modifies the Config struct.
type ConfigOption func(*Config)

// WithJokersInDeck sets the number of jokers in a deck.
func WithJokersInDeck(jokerCount int) ConfigOption {
	return func(c *Config) {
		c.jokerInDeckCount = jokerCount
	}
}

// WithCompareFunc sets the comparison function for sorting the deck.
// If shuffling is enabled as well, the deck will be sorted after shuffling.
func WithCompareFunc(cmp func(a, b Card) int) ConfigOption {
	return func(c *Config) {
		c.cmpFunc = cmp
	}
}

// WithShuffle sets whether to shuffle the deck.
func WithShuffle(shuffle bool) ConfigOption {
	return func(c *Config) {
		c.shuffleDeck = shuffle
	}
}

// WithExclude sets the function to exclude cards from the deck.
func WithExclude(exclude func(Card) bool) ConfigOption {
	return func(c *Config) {
		c.excludeFunc = exclude
	}
}

// WithAdditionalDecks sets the number of additional decks to create.
func WithAdditionalDecks(additionalDecksCount int) ConfigOption {
	return func(c *Config) {
		c.additionalDecksCount = additionalDecksCount
	}
}

// NewDeck creates a new deck with the given options.
func NewDeck(opts ...ConfigOption) *Deck {
	cfg := Config{}
	for _, o := range opts {
		o(&cfg)
	}

	deck := buildDecks(cfg)

	deck = modifyDeckContent(cfg, deck)

	deck = modifyDeckOrder(cfg, deck)

	return deck
}

// buildDecks creates a new deck with the given configuration.
// It creates a base deck and adds additional decks if specified.
func buildDecks(c Config) *Deck {
	deck := newDeckWithCards(c)

	for range c.additionalDecksCount {
		addDeck := newDeckWithCards(c)
		deck.Cards = append(deck.Cards, addDeck.Cards...)
	}

	return deck
}

// newDeckWithCards creates a new deck with cards based on the configuration.
func newDeckWithCards(c Config) *Deck {
	const numberOfCardsInDeck = 13

	deck := newEmptyDeck()

	// Create numberOfCardsInDeck(13) cards for each suit.
	for _, suit := range []Suit{Hearts, Spades, Diamonds, Clubs} {
		for value := range numberOfCardsInDeck {
			// The loop will create cards with values from 0 to 12, but we want 1-13, so we
			// increment by 1.
			deck.Cards = append(deck.Cards, NewCard(suit, uint8(value)+1))
		}
	}

	// Create as many jokers as specified in the configuration.
	// By default, there are no jokers in a deck
	for range c.jokerInDeckCount {
		deck.Cards = append(deck.Cards, NewCard(Jokers, 1))
	}

	return deck
}

// newEmptyDeck creates a new empty deck.
func newEmptyDeck() *Deck {
	return &Deck{
		Cards: make([]Card, 0),
	}
}

// modifyDeckContent modifies the deck content based on the configuration.
func modifyDeckContent(cfg Config, deck *Deck) *Deck {
	if cfg.excludeFunc != nil {
		for cardIndex, card := range deck.Cards {
			if cfg.excludeFunc(card) {
				deck.Cards = append(deck.Cards[:cardIndex], deck.Cards[cardIndex+1:]...)
			}
		}
	}

	return deck
}

// modifyDeckOrder modifies the order of the deck based on the configuration. It
// should be called after the contents of the deck are finalized.
func modifyDeckOrder(cfg Config, deck *Deck) *Deck {
	if cfg.shuffleDeck {
		rand.Shuffle(len(deck.Cards), func(i, j int) { deck.Cards[i], deck.Cards[j] = deck.Cards[j], deck.Cards[i] })
	}

	if cfg.cmpFunc != nil {
		slices.SortStableFunc(deck.Cards, cfg.cmpFunc)
	}

	return deck
}

// FetchCard fetches the top card from the deck, once a card has been fetched it
// is no longer in the deck.
//
// If the deck is empty, it returns an empty card and false.
func (d *Deck) FetchCard() (Card, bool) {
	var card Card
	if len(d.Cards) == 0 {
		return Card{}, false
	}
	if len(d.Cards) == 1 {
		card = d.Cards[0]
		d.Cards = make([]Card, 0)
		return card, true
	}

	// pull the card from the end of the deck,
	card = d.Cards[len(d.Cards)-1]
	// And remove the last card from the deck.
	d.Cards = d.Cards[:len(d.Cards)-1]

	return card, true
}

// PushCard pushes a card to the bottom of the deck.
func (d *Deck) PushCard(card Card) {
	d.Cards = append([]Card{card}, d.Cards...)
}
