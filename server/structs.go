package main

type Card struct {
	Suit string `json:"suit"`
	Rank string `json:"rank"`
}

type ReceivedCards struct {
	Catch bool   `json:"catch"`
	Cards []Card `json:"cards"`
}

type Hand struct {
	Cards []Card `json:"cards"`
}

type Game struct {
	Id      string   `json:"game_id"`
	Players []string `json:"players"`
}

type CardRequest struct {
	OtherPlayer string `json:"player"`
	Rank        string `json:"rank"`
}

// map deck of cards "values" to our "rank" spec
var valueMap = map[string]string{
	"ACE":   "ace",
	"2":     "two",
	"3":     "three",
	"4":     "four",
	"5":     "five",
	"6":     "six",
	"7":     "seven",
	"8":     "eight",
	"9":     "nine",
	"0":     "ten",
	"JACK":  "jack",
	"QUEEN": "queen",
	"KING":  "king",
}

// map our "rank" spec to deck of card "values"
var rankMap = map[string]string{
	"ace":   "ACE",
	"two":   "2",
	"three": "3",
	"four":  "4",
	"five":  "5",
	"six":   "6",
	"seven": "7",
	"eight": "8",
	"nine":  "9",
	"ten":   "0",
	"jack":  "JACK",
	"queen": "QUEEN",
	"king":  "KING",
}
