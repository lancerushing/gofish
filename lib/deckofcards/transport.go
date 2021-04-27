package deckofcards

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog/log"
)

type LoggingRoundTripper struct {
	Next http.RoundTripper
}

func (lrt LoggingRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, e error) {
	// Do "before sending requests" actions here.
	log.Debug().Interface("url", req.URL).Str("method", req.Method).Msg("sending")

	resp, e = lrt.Next.RoundTrip(req)

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {
		log.Error().Bytes("response body", b).Send()
	} else {
		log.Debug().Bytes("response body", b).Send()
	}

	resp.Body = ioutil.NopCloser(bytes.NewReader(b))

	return
}
