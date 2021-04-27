package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/log"
	"github.com/lancerushing/gofish/lib"
	"github.com/lancerushing/gofish/lib/deckofcards"
)

func main() {
	lib.LogSetup()
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Run Error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Info().Msgf("Defaulting to port: %s", port)
	}

	mux := chi.NewRouter()
	mux.Use(middleware.Logger)

	client := deckofcards.MakeClient("https://deckofcardsapi.com/")

	s := server{
		client: &client,
	}

	mux.Post("/games", s.handleNewGame())
	mux.Get("/games/{gameId}/players/{playerName}", s.handlePlayerHand())
	mux.Post("/games/{gameId}/players/{playerName}/fish", s.handleFish())

	log.Info().Msgf("Listening on port %s", port)

	return http.ListenAndServe(":"+port, mux)
}
