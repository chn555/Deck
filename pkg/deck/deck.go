package deck

import (
	"math/rand"
)

// Deck represents a deck of cards.
type Deck struct {
	// Cards represents the cards in the deck.
	Cards []Card
}

// Config represents the configuration for the deck.
type Config struct {
	// AdditionalDecks is the number of decks to create in addition to the base
	// deck. A base deck is always created, and when AdditionalDecks is 0, only
	// the base deck is created. If AdditionalDecks is not 0, additional decks
	// are created and added to the base deck.
	AdditionalDecks int
	// JokerInDeckCount is the number of jokers a single deck, by default there are
	// no jokers in a deck. If AdditionalDecks is used, this will be the numbers
	// of jokers in each additional deck as well.
	JokerInDeckCount int
	// ExcludeCardFunc is the function used to exclude cards from the deck.
	ExcludeCardFunc func(Card) bool
	// ShuffleDeck indicates whether to shuffle the deck.
	// If another option was provided that modifies the order of the deck,
	// the deck will be shuffled before that option is applied.
	// If AdditionalDecks is not 0, the deck will be shuffled after
	// the additional decks are added.
	ShuffleDeck bool
}

// ConfigOption is a function that modifies the Config struct.
type ConfigOption func(*Config)

// WithAdditionalDecks is the number of decks to create in addition to the base
// deck. A base deck is always created, and when AdditionalDecks is 0, only
// the base deck is created. If AdditionalDecks is not 0, additional decks
// are created and added to the base deck.
func WithAdditionalDecks(additionalDecks int) ConfigOption {
	return func(c *Config) {
		c.AdditionalDecks = additionalDecks
	}
}

// WithJokersInDeck sets the number of jokers a single deck, by default there are
// no jokers in a deck. If AdditionalDecks is used, this will be the numbers
// of jokers in each additional deck as well.
func WithJokersInDeck(jokerCount int) ConfigOption {
	return func(c *Config) {
		c.JokerInDeckCount = jokerCount
	}
}

// WithExclude sets the function to exclude cards from the deck.
func WithExclude(exclude func(Card) bool) ConfigOption {
	return func(c *Config) {
		c.ExcludeCardFunc = exclude
	}
}

// WithShuffle indicates whether to shuffle the deck.
// If another option was provided that modifies the order of the deck,
// the deck will be shuffled before that option is applied.
// If AdditionalDecks is not 0, the deck will be shuffled after
// the additional decks are added.
func WithShuffle(shuffle bool) ConfigOption {
	return func(c *Config) {
		c.ShuffleDeck = shuffle
	}
}

// NewDeck creates a new deck with the given options.
func NewDeck(opts ...ConfigOption) *Deck {
	cfg := Config{}
	for _, o := range opts {
		o(&cfg)
	}

	// Create a new deck with the given configuration.
	// The order of the operations is important:
	// 1. Create the base deck.
	// 2. Add additional decks as needed.
	// 3. Modify the contents of the deck as needed, adding jokers, excluding cards, etc.
	// 4. Modify the order of the deck as needed, shuffling or sorting. Deterministic sorting should come last.
	deck := newDeckWithCards().
		addAdditionalDecks(cfg).
		addJokers(cfg).
		excludeCards(cfg).
		shuffleDeck(cfg)

	return deck
}

// The function newDeckWithCards creates a new deck with default cards.
//
// The default deck is a standard 52 card deck with 4 suits and 13 values.
func newDeckWithCards() *Deck {
	// The numbers of cards of each suit in a standard French-suited deck.
	// https://en.wikipedia.org/wiki/Standard_52-card_deck
	const numberOfCardsPerSuite = 13

	deck := newEmptyDeck()

	// Create numberOfCardsPerSuite(13) cards for each suit.
	for _, suit := range []Suit{Hearts, Spades, Diamonds, Clubs} {
		for value := range numberOfCardsPerSuite {
			// The loop will create cards with values from 0 to 12, but we want 1-13, so we
			// increment by 1.
			deck.Cards = append(deck.Cards, NewCard(suit, uint8(value)+1))
		}
	}

	return deck
}

// The function newEmptyDeck creates a new empty deck.
// It provides a valid, initialized deck, but without any cards.
func newEmptyDeck() *Deck {
	return &Deck{
		Cards: make([]Card, 0),
	}
}

// The function addAdditionalDecks creates additional decks as needed. The decks
// it creates are default decks, with 52 cards each, which are then appended to
// the base deck.
func (d *Deck) addAdditionalDecks(c Config) *Deck {
	for range c.AdditionalDecks {
		addDeck := newDeckWithCards()
		d.Cards = append(d.Cards, addDeck.Cards...)
	}

	return d
}

func (d *Deck) addJokers(cfg Config) *Deck {
	// We need to create as many jokers as specified in the configuration, for each
	// deck, so we have the base deck plus however many additional decks are
	// configured.
	totalNumberOfDecks := 1 + cfg.AdditionalDecks

	// Multiplying the configured number of jokers by the total number of decks.
	totalNumberOfJokers := cfg.JokerInDeckCount * totalNumberOfDecks

	// Create as many jokers as we determined are needed.
	// By default, the decks are created without jokers.
	for range totalNumberOfJokers {
		d.Cards = append(d.Cards, NewCard(Jokers, 1))
	}

	return d
}

// The function excludeCards removes cards if WithExclude is set.
// It should be called before the order is modified.
func (d *Deck) excludeCards(cfg Config) *Deck {
	if cfg.ExcludeCardFunc == nil {
		return d
	}

	for cardIndex, card := range d.Cards {
		if cfg.ExcludeCardFunc(card) {
			// To remove the card, we need to create a new slice without the card. This is
			// done by appending the slice before the card and the slice after the card.
			d.Cards = append(d.Cards[:cardIndex], d.Cards[cardIndex+1:]...)
		}
	}

	return d
}

// If WithShuffle is set, the deck is shuffled.
// This function should be called after the contents of the deck are finalized.
func (d *Deck) shuffleDeck(cfg Config) *Deck {
	if cfg.ShuffleDeck {
		// Using rand.Shuffle here, we provide a function that swaps the cards in the
		// deck slice
		rand.Shuffle(len(d.Cards),
			func(i, j int) {
				d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
			},
		)
	}

	return d
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
