package deck

// Deck represents a deck of cards.
type Deck struct {
	// Cards represents the cards in the deck.
	Cards []Card
}

// Config represents the configuration for the deck.
type Config struct{}

// NewDeck creates a new deck with the given options.
func NewDeck() *Deck {
	deck := newDeckWithCards()

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
