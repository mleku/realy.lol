// Package normalize is a set of tools for cleaning up URL s and formatting
// nostr OK and CLOSED messages.
package normalize

import (
	"bytes"
	"net/url"

	"orly.dev/chk"
	"orly.dev/ints"
	"orly.dev/log"
)

var (
	hp    = bytes.HasPrefix
	WS    = []byte("ws://")
	WSS   = []byte("wss://")
	HTTP  = []byte("http://")
	HTTPS = []byte("https://")
)

// URL normalizes the URL
//
// - Adds wss:// to addresses without a port, or with 443 that have no protocol prefix
//
// - Adds ws:// to addresses with any other port
//
// - Converts http/s to ws/s
func URL[V string | []byte](v V) (b []byte) {
	u := []byte(v)
	if len(u) == 0 {
		return nil
	}
	u = bytes.TrimSpace(u)
	u = bytes.ToLower(u)
	// if address has a port number, we can probably assume it is insecure websocket as most
	// public or production relays have a domain name and a well known port 80 or 443 and thus
	// no port number.
	//
	// if a protocol prefix is present, we assume it is already complete. Converting http/s to
	// websocket equivalent will be done later anyway.
	if bytes.Contains(u, []byte(":")) &&
		!(hp(u, HTTP) || hp(u, HTTPS) || hp(u, WS) || hp(u, WSS)) {

		split := bytes.Split(u, []byte(":"))
		if len(split) != 2 {
			log.D.F("Error: more than one ':' in URL: '%s'", u)
			// this is a malformed URL if it has more than one ":", return empty
			// since this function does not return an error explicitly.
			return
		}
		p := ints.New(0)
		_, err := p.Unmarshal(split[1])
		if chk.E(err) {
			log.D.F("Error normalizing URL '%s': %s", u, err)
			// again, without an error we must return nil
			return
		}
		if p.Uint64() > 65535 {
			log.D.F(
				"Port on address %d: greater than maximum 65535",
				p.Uint64(),
			)
			return
		}
		// if the port is explicitly set to 443 we assume it is wss:// and drop the port.
		if p.Uint16() == 443 {
			u = append(WSS, split[0]...)
		} else {
			u = append(WSS, u...)
		}
	}

	// if prefix isn't specified as http/s or websocket, assume secure websocket and add wss
	// prefix (this is the most common).
	if !(hp(u, HTTP) || hp(u, HTTPS) || hp(u, WS) || hp(u, WSS)) {
		u = append(WSS, u...)
	}
	var err error
	var p *url.URL
	if p, err = url.Parse(string(u)); chk.E(err) {
		return
	}
	// convert http/s to ws/s
	switch p.Scheme {
	case "https":
		p.Scheme = "wss"
	case "http":
		p.Scheme = "ws"
	}
	// remove trailing path slash
	p.Path = string(bytes.TrimRight([]byte(p.Path), "/"))
	return []byte(p.String())
}
