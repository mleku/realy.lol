// Copyright (c) 2015-2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package btcec

import (
	"bytes"
	"testing"

	"orly.dev/chk"
)

func TestGenerateSharedSecret(t *testing.T) {
	privKey1, err := NewSecretKey()
	if chk.E(err) {
		t.Errorf("secret key generation error: %s", err)
		return
	}
	privKey2, err := NewSecretKey()
	if chk.E(err) {
		t.Errorf("secret key generation error: %s", err)
		return
	}
	secret1 := GenerateSharedSecret(privKey1, privKey2.PubKey())
	secret2 := GenerateSharedSecret(privKey2, privKey1.PubKey())
	if !bytes.Equal(secret1, secret2) {
		t.Errorf(
			"ECDH failed, secrets mismatch - first: %x, second: %x",
			secret1, secret2,
		)
	}
}
