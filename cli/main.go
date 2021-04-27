package main

// This program will play a complete game of 'go fish'
// go run ./cli/ new Brian,Nicholas,Lance

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
	"github.com/lancerushing/gofish/lib"
	"github.com/lancerushing/gofish/lib/deckofcards"
)

const sep string = "############################"

func main() {
	lib.LogSetup()

	args := os.Args
	if len(args) != 3 {
		log.Fatal().Int("length", len(args)).Msg("invalid arg length")
	}

	// if the deckofcards.com server crashes, you can re-use the previous ID to continue
	// use "new" to start a fresh game
	deckID := args[1]

	playerNames := strings.Split(args[2], ",")

	startingHandSize := 5 // should be odd number, prevent occasional instant win

	if err := run(deckID, playerNames, startingHandSize); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Run Error: %s\n", err)
		os.Exit(1)
	}
}

func run(deckID string, playerNames []string, startingHandSize int) error {
	c := deckofcards.MakeClient("https://deckofcardsapi.com/")
	deck := setupDeck(c, deckID, playerNames, startingHandSize)

	winnerName := gameLoop(c, deck, playerNames)
	fmt.Printf("%v %s\n", color.GreenString("Winner:"), winnerName)

	return nil
}

func gameLoop(c deckofcards.Client, deck deckofcards.DeckResponse, playerNames []string) (winnerName string) {

	winnerName = checkHands(c, deck, playerNames)
	if winnerName != "" {
		return winnerName
	}

	for  {
		for _, currentPlayerName := range playerNames {

			winnerName = checkHands(c, deck, playerNames)
			if winnerName != "" {
				return winnerName
			}

			fmt.Printf("%s\n%s's turn\n", sep, currentPlayerName)
			time.Sleep(1 * time.Second)

			var foundPair bool
			for ok := true; ok; ok = foundPair { // keep asking if player finds a pair
				foundPair = false

				currentPlayerCards := c.ListPile2(deck, currentPlayerName).Cards

				// pick "other" player
				otherPlayer := getOther(currentPlayerName, playerNames)
				card := pickCard(currentPlayerName, currentPlayerCards)

				fmt.Printf("%s asks %s for: %s\n", currentPlayerName, otherPlayer, card.Value)
				time.Sleep(1 * time.Second)

				var foundCard deckofcards.Card
				otherPlayerCards := c.ListPile2(deck, otherPlayer).Cards
				foundCard, foundPair = deckofcards.PileContainsValue(otherPlayerCards, card.Value)
				if foundPair {
					fmt.Printf("%s has a %s\n", otherPlayer, card.Value)
					c.DrawFromPile(deck, currentPlayerName, card.Code) // discard the card
					c.DrawFromPile(deck, otherPlayer, foundCard.Code)  // discard the card

					// check for winner
					if len(currentPlayerCards) == 1 { // assume if list was a length of 1, it is now 0 (prevent extra api call)
						return currentPlayerName
					}

					if len(otherPlayerCards) == 1 { // assume if list was a length of 1, it is now 0 (prevent extra api call)
						return otherPlayer
					}

				} else {
					fmt.Printf("%s says, \"Go fish!\"\n", otherPlayer)
					draw := c.Draw(deck, 1)

					fmt.Printf("%s draws a %s\n", currentPlayerName, draw.Cards[0].Code)
					c.AddToPile(deck, currentPlayerName, draw.Cards)
				}
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func setupDeck(c deckofcards.Client, deckId string, players []string, numOfStartingCard int) deckofcards.DeckResponse {
	if deckId != "new" {
		return deckofcards.DeckResponse{
			ID: deckId,
		}
	}

	fmt.Print("shuffling\n")
	sfl := c.Shuffle(deckId)
	deck := sfl.DeckResponse
	color.Red("deck id = %s", deck.ID)

	fmt.Print("dealing cards\n")
	for _, player := range players {
		draw := c.Draw(deck, numOfStartingCard)
		c.AddToPile(deck, player, draw.Cards)

		playerCards := c.ListPile2(deck, player).Cards
		printCards(player, playerCards)
	}

	return deck
}

// keep track of next card to pick
var nextIndexToPick map[string]int

// pickCard from the pile
func pickCard(player string, cards []deckofcards.Card) deckofcards.Card {
	if nextIndexToPick == nil {
		nextIndexToPick = map[string]int{}
	}

	i := nextIndexToPick[player]

	if i >= len(cards) {
		i = 0
	}

	nextIndexToPick[player] = i + 1
	return cards[i]

	// we could pick a random card from cards, but that makes for very long games
	// s := rand.NewSource(time.Now().Unix())
	// r := rand.New(s)
	// i := r.Intn(len(cards))
	// return cards[i]
}

func getOther(currentPlayer string, players []string) string {
	var others []string

	for _, player := range players {
		if currentPlayer != player {
			others = append(others, player)
		}
	}

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	i := r.Intn(len(others))

	return others[i]
}

// checkHands will check the hands of each player, and remove any pairs
func checkHands(c deckofcards.Client, deck deckofcards.DeckResponse, playerNames []string) (winner string) {
	for _, playerName := range playerNames {
		var foundPair bool

		for ok := true; ok; ok = foundPair {
			foundPair = false

			playerCards := c.ListPile2(deck, playerName).Cards
			printCards(playerName, playerCards)

			if len(playerCards) == 0 {
				return playerName
			}

			foundPair = discardMatches(c, deck, playerName, playerCards)
		}
	}

	return ""
}

func discardMatches(c deckofcards.Client, deck deckofcards.DeckResponse, playerName string, playerCards []deckofcards.Card) (foundPair bool) {
	for i := 0; i < len(playerCards)-1; i++ {
		c1 := playerCards[i]
		for ii := i + 1; ii < len(playerCards); ii++ {
			c2 := playerCards[ii]
			if c1.Value == c2.Value {
				fmt.Printf("    Match Found: %s\n", c1.Value)
				c.DrawFromPile(deck, playerName, c1.Code)
				c.DrawFromPile(deck, playerName, c2.Code)
				return true
			}
		}
	}
	return false
}

func printCards(player string, cards []deckofcards.Card) {
	var codes []string
	for _, card := range cards {
		codes = append(codes, card.Code)
	}

	fmt.Printf("    %-10s %+v\n", player, codes)
}
