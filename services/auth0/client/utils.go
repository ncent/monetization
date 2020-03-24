package client

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
)

func GeneratePrivateAndPubKeyHex() (string, string, error) {
	pubkeyCurve := elliptic.P256()

	privatekey, err := ecdsa.GenerateKey(pubkeyCurve, rand.Reader)

	if err != nil {
		return "", "", err
	}

	var pubkey ecdsa.PublicKey
	pubkey = privatekey.PublicKey

	pk := hex.EncodeToString(crypto.FromECDSA(privatekey))
	pa := hex.EncodeToString(crypto.FromECDSAPub(&pubkey))

	return pk, pa, nil
}
