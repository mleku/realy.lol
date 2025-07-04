package btcec_test

import (
	"bytes"
	"testing"
	"time"

	"orly.dev/chk"
	"orly.dev/log"
	"orly.dev/p256k/btcec"
)

func TestSigner_Generate(t *testing.T) {
	for _ = range 100 {
		var err error
		signer := &btcec.Signer{}
		var skb []byte
		if err = signer.Generate(); chk.E(err) {
			t.Fatal(err)
		}
		skb = signer.Sec()
		if err = signer.InitSec(skb); chk.E(err) {
			t.Fatal(err)
		}
	}
}

func TestBTCECSignerVerify(t *testing.T) {
}

func TestBTCECSignerSign(t *testing.T) {
}

func TestBTCECECDH(t *testing.T) {
	n := time.Now()
	var err error
	var counter int
	const total = 200
	var count int
	for _ = range total {
		s1 := new(btcec.Signer)
		if err = s1.Generate(); chk.E(err) {
			t.Fatal(err)
		}
		s2 := new(btcec.Signer)
		if err = s2.Generate(); chk.E(err) {
			t.Fatal(err)
		}
		for _ = range total {
			var secret1, secret2 []byte
			if secret1, err = s1.ECDH(s2.Pub()); chk.E(err) {
				t.Fatal(err)
			}
			if secret2, err = s2.ECDH(s1.Pub()); chk.E(err) {
				t.Fatal(err)
			}
			if !bytes.Equal(secret1, secret2) {
				counter++
				t.Errorf(
					"ECDH generation failed to work in both directions, %x %x",
					secret1,
					secret2,
				)
			}
			count++
		}
	}
	a := time.Now()
	duration := a.Sub(n)
	log.I.Ln(
		"errors", counter,
		"total", count,
		"time", duration,
		"time/op", duration/time.Duration(count),
		"ops/sec", int(time.Second)/int(duration/time.Duration(count)),
	)
}
