package deck

// Deck represents a deck of cards.
type Deck struct {
	// Cards represents the cards in the deck.
	Cards []Card
}

// FetchCard fetches the top card from the deck,
// once a card has been fetched it is no longer in the deck.
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
