package deck

import (
	"math/rand"
	"slices"
)

// Deck represents a deck of cards
type Deck struct {
	Cards []Card
}

// FetchCard fetches the top card from the deck
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

	card, d.Cards = d.Cards[len(d.Cards)-1], d.Cards[:len(d.Cards)-1]
	return card, true
}

// PushCard pushes a card to the start of the deck
func (d *Deck) PushCard(card Card) {
	d.Cards = append([]Card{card}, d.Cards...)
}

// Config represents the configuration for the deck
type Config struct {
	cmpFunc              func(a, b Card) int
	shuffleDeck          bool
	jokerInDeckCount     int
	excludeFunc          func(Card) bool
	additionalDecksCount int
}

type ConfigOption func(*Config)

// WithJokersInDeck sets the number of jokers in a deck
func WithJokersInDeck(jokerCount int) ConfigOption {
	return func(c *Config) {
		c.jokerInDeckCount = jokerCount
	}
}

// WithCompareFunc sets the comparison function for sorting the deck
// If shuffling is enabled as well, the deck will be sorted after shuffling
func WithCompareFunc(cmp func(a, b Card) int) ConfigOption {
	return func(c *Config) {
		c.cmpFunc = cmp
	}
}

// WithShuffle sets whether to shuffle the deck
func WithShuffle(shuffle bool) ConfigOption {
	return func(c *Config) {
		c.shuffleDeck = shuffle
	}
}

// WithExclude sets the function to exclude cards from the deck
func WithExclude(exclude func(Card) bool) ConfigOption {
	return func(c *Config) {
		c.excludeFunc = exclude
	}
}

// WithAdditionalDecks sets the number of additional decks to create
func WithAdditionalDecks(additionalDecksCount int) ConfigOption {
	return func(c *Config) {
		c.additionalDecksCount = additionalDecksCount
	}
}

// NewDeck creates a new deck with the given options
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

func buildDecks(c Config) *Deck {
	deck := newDeck(c)

	for range c.additionalDecksCount {
		addDeck := newDeck(c)
		deck.Cards = append(deck.Cards, addDeck.Cards...)
	}

	return deck
}

func newDeck(c Config) *Deck {
	deck := newEmptyDeck()

	for _, suit := range []Suit{Aces, Spades, Diamonds, Clubs} {
		for value := range 13 {
			deck.Cards = append(deck.Cards, NewCard(suit, uint8(value)+1))
		}
	}

	for range c.jokerInDeckCount {
		deck.Cards = append(deck.Cards, NewCard(Jokers, 1))
	}

	return deck
}

func newEmptyDeck() *Deck {
	return &Deck{
		Cards: make([]Card, 0),
	}
}

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

func modifyDeckOrder(cfg Config, deck *Deck) *Deck {
	if cfg.shuffleDeck {
		rand.Shuffle(len(deck.Cards), func(i, j int) { deck.Cards[i], deck.Cards[j] = deck.Cards[j], deck.Cards[i] })
	}

	if cfg.cmpFunc != nil {
		slices.SortStableFunc(deck.Cards, cfg.cmpFunc)
	}

	return deck
}
