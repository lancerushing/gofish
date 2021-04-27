package deckofcards

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
	"github.com/lancerushing/gofish/lib"
)

func TestClientIntegration(t *testing.T) {
	lib.LogSetup()

	should := is.New(t)
	c := MakeClient("https://deckofcardsapi.com/")

	shuffle := c.Shuffle("")

	deck := shuffle.DeckResponse

	t.Logf("%+v\n", shuffle)
	should.Equal(true, shuffle.Success)

	draw := c.Draw(deck, 2)
	t.Logf("%+v\n", draw)
	should.Equal(2, len(draw.Cards))

	resp := c.AddToPile(deck, "lance", draw.Cards)
	t.Logf("%+v\n", resp)

	pile := c.ListPile(deck, "lance")
	t.Logf("%+v\n", pile)
}

func TestClient(t *testing.T) {
	should := is.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log(r.URL.Path)
		t.Log(r.URL.RawQuery)
		fmt.Fprintln(w, "{}")
	}))
	defer ts.Close()

	c := MakeClient(ts.URL)

	deck := DeckResponse{
		ID: "foo",
	}

	draw := c.Draw(deck, 2)
	t.Logf("%+v\n", draw)
	should.Equal(2, len(draw.Cards))
}
