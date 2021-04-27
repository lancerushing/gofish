package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/lancerushing/gofish/lib"
	"github.com/lancerushing/gofish/lib/deckofcards"
)

type server struct {
	client *deckofcards.Client
}

type gameSettings struct {
	Players []string
}

func (s *server) handleNewGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// load in request
		var newGameSettings gameSettings
		err := json.NewDecoder(r.Body).Decode(&newGameSettings)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// validate
		if !(len(newGameSettings.Players) >= 2 && len(newGameSettings.Players) <= 5) {
			http.Error(w, "players must be between 2 and 5", http.StatusBadRequest)
			return
		}
		if !lib.Unique(newGameSettings.Players) {
			http.Error(w, "player names must be unique", http.StatusBadRequest)
			return
		}

		// create a new deck from deckofcars.com
		deckId := "new"
		sfl := s.client.Shuffle(deckId)
		deck := sfl.DeckResponse

		game := Game{
			Id:      sfl.ID,
			Players: newGameSettings.Players,
		}

		for _, player := range newGameSettings.Players {
			draw := s.client.Draw(deck, 5)
			s.client.AddToPile(deck, player, draw.Cards)
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(game)
	}
}

func (s *server) handlePlayerHand() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameId := chi.URLParam(r, "gameId")
		playerName := chi.URLParam(r, "playerName")

		deck := deckofcards.DeckResponse{
			ID: gameId,
		}
		playerCards := s.client.ListPile2(deck, playerName).Cards

		var cards Hand
		for _, playerCard := range playerCards {
			card := Card{
				Suit: strings.ToLower(playerCard.Suit),
				Rank: valueMap[playerCard.Value],
			}
			cards.Cards = append(cards.Cards, card)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cards)
	}
}

func (s *server) handleFish() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// load in request
		var cardReq CardRequest
		err := json.NewDecoder(r.Body).Decode(&cardReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// validate
		if cardReq.OtherPlayer == "" {
			http.Error(w, "player name is empty", http.StatusBadRequest)
			return
		}

		value, ok := rankMap[cardReq.Rank]
		if !ok {
			http.Error(w, "rank is empty", http.StatusBadRequest)
			return
		}

		gameId := chi.URLParam(r, "gameId")
		playerName := chi.URLParam(r, "playerName")

		deck := deckofcards.DeckResponse{
			ID: gameId,
		}
		otherCards := s.client.ListPile2(deck, cardReq.OtherPlayer).Cards

		var rc ReceivedCards

		foundCard, found := deckofcards.PileContainsValue(otherCards, value)
		if !found {
			w.WriteHeader(http.StatusCreated) // maybe change to 200, 202, or 404
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(rc)
			return
		}

		rc.Catch = true
		rc.Cards = append(rc.Cards, Card{
			Suit: strings.ToLower(foundCard.Suit),
			Rank: valueMap[foundCard.Value],
		})

		s.client.AddToPile(deck, playerName, []deckofcards.Card{foundCard})

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rc)
	}
}
