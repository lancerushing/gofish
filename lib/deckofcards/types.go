package deckofcards

type DeckResponse struct {
	Success   bool   `json:"success"`
	ID        string `json:"deck_id"`
	Remaining int    `json:"remaining"`
}

type ShuffleResponse struct {
	DeckResponse
	Shuffled bool `json:"shuffled"`
}

type DrawResponse struct {
	DeckResponse
	Cards []Card `json:"cards"`
}

type PileResponse struct {
	DeckResponse
	Piles struct {
		Discard struct {
			Remaining int
		}
	}
}

type CardCollection struct {
	Remaining int
	Cards     []Card
}

type PilesResponse struct {
	DeckResponse
	Piles map[string]CardCollection
}

type PilesDiscardResponse struct {
	DeckResponse
	Cards []Card
}

type Card struct {
	Image string `json:"image"`
	Value string `json:"value"`
	Suit  string `json:"suit"`
	Code  string `json:"code"`
}

func PileContainsValue(cards []Card, value string) (card Card, found bool) {
	for _, card = range cards {
		if card.Value == value {
			return card, true
		}
	}
	return
}
