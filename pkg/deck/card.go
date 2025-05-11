package deck

type Suit int8

const (
	UnknownSuit Suit = iota
	Aces
	Spades
	Diamonds
	Clubs
	Jokers
)

// Card represents a playing card
type Card struct {
	// Suit represents the card suit
	Suit Suit
	// Values represents the card value
	// Jacks are 11, Queens are 12, Kings are 13
	// For cards of suit Jokers the value does not matter
	Value uint8
}

// NewCard creates a new card with the given suit and value
func NewCard(suit Suit, value uint8) Card {
	return Card{
		Suit:  suit,
		Value: value,
	}
}
