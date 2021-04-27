package deckofcards

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/lancerushing/gofish/lib"
)

type Client struct {
	BaseURL   *url.URL
	UserAgent string

	httpClient *http.Client
}

func MakeClient(base string) Client {
	b, err := url.Parse(base)
	if err != nil {
		panic(err) // not for production
	}

	c := Client{
		BaseURL:    b,
		UserAgent:  "vericred",
		httpClient: &http.Client{
			// Transport: LoggingRoundTripper{http.DefaultTransport},
		},
	}

	return c
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	return req, nil
}

const maxTries = 5

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := int64(1); i <= maxTries; i++ { // the api server fails often.
		time.Sleep(time.Duration(i*100) * time.Millisecond) // limit speed, and backoff after server serror
		resp, err = c.doOnce(req, v)
		if err == nil {
			if resp.StatusCode != http.StatusInternalServerError {
				break
			}
		}
	}

	return resp, err
}

func (c *Client) doOnce(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	lib.CheckError(err)

	err = json.Unmarshal(b, v)
	if err != nil {
		log.Error().Bytes("json", b).Msg("json unmarshal error")
	}

	return resp, err
}

// Shuffle will shuffle dec
func (c *Client) Shuffle(deckID string) ShuffleResponse {
	var resp ShuffleResponse

	path := fmt.Sprintf("/api/deck/%s/shuffle/", url.PathEscape(deckID))

	req, _ := c.newRequest("GET", path, nil)
	_, err := c.do(req, &resp) // TODO add error handling
	lib.CheckError(err)

	return resp
}

// Shuffle will shuffle dec
func (c *Client) Draw(deck DeckResponse, count int) DrawResponse {
	resp := DrawResponse{}

	path := fmt.Sprintf("/api/deck/%s/draw/", url.PathEscape(deck.ID))

	req, _ := c.newRequest("GET", path, nil)
	q := req.URL.Query()
	q.Add("count", strconv.Itoa(count))
	req.URL.RawQuery = q.Encode()

	_, err := c.do(req, &resp)
	lib.CheckError(err)

	return resp
}

// AddToPile ...
func (c *Client) AddToPile(deck DeckResponse, pileName string, cards []Card) PileResponse {
	resp := PileResponse{}

	path := fmt.Sprintf("/api/deck/%s/pile/%s/add/", url.PathEscape(deck.ID), url.PathEscape(pileName))
	req, _ := c.newRequest("GET", path, nil)

	var codes []string

	for _, card := range cards {
		codes = append(codes, card.Code)
	}

	q := req.URL.Query()
	q.Add("cards", strings.Join(codes, ","))

	req.URL.RawQuery = q.Encode()

	_, _ = c.do(req, &resp) // TODO add error handling
	return resp
}

// ListPile ...
func (c *Client) ListPile(deck DeckResponse, pileName string) PilesResponse {
	resp := PilesResponse{}

	path := fmt.Sprintf("/api/deck/%s/pile/%s/list/", url.PathEscape(deck.ID), url.PathEscape(pileName))
	req, _ := c.newRequest("GET", path, nil)

	_, _ = c.do(req, &resp) // TODO add error handling
	return resp
}

// ListPile2 ...
func (c *Client) ListPile2(deck DeckResponse, pileName string) CardCollection {
	resp := c.ListPile(deck, pileName)
	return resp.Piles[pileName]
}

// DrawFromPile ...
func (c *Client) DrawFromPile(deck DeckResponse, pileName string, cardCode string) PilesDiscardResponse {
	resp := PilesDiscardResponse{}

	path := fmt.Sprintf("/api/deck/%s/pile/%s/draw/", url.PathEscape(deck.ID), url.PathEscape(pileName))
	req, _ := c.newRequest("GET", path, nil)
	q := req.URL.Query()
	q.Add("cards", cardCode)
	req.URL.RawQuery = q.Encode()

	_, _ = c.do(req, &resp) // TODO add error handling
	return resp
}
