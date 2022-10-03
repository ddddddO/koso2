package koso2

import (
	"bufio"
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/crypto/ssh"
)

// NOTE: 一旦、RSAのみ考える
func Run(ghUserID, plainMessage string, callback func(encrypted string) error) error {
	pubKeys, err := fetchPublicKeys(ghUserID)
	if err != nil {
		return err
	}

	pk := pubKeys[0] // ローカルにあるのがこれなので、一旦
	rsaPubKey, err := parsePublicKey(pk)
	if err != nil {
		return err
	}

	encrypted, err := encryptMessage(plainMessage, rsaPubKey)
	if err != nil {
		return err
	}

	if err := callback(encrypted); err != nil {
		return err
	}

	return nil
}

func fetchPublicKeys(ghUserID string) ([]string, error) {
	u := &url.URL{
		Scheme: "https",
		Host:   "github.com",
		Path:   fmt.Sprintf("%s.keys", ghUserID),
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to request. status code: %d\n", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	pubKeys := []string{}
	for scanner.Scan() {
		pubKeys = append(pubKeys, scanner.Text())
	}

	return pubKeys, nil
}

func parsePublicKey(raw string) (crypto.PublicKey, error) {
	publicKey := strings.TrimLeft(raw, "ssh-rsa ")
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(publicKey)))
	n, err := base64.StdEncoding.Decode(decoded, []byte(publicKey))
	if err != nil {
		return nil, fmt.Errorf("failed to base64 decode.\n%v", err)
	}
	decoded = decoded[:n]

	sshPubKey, err := ssh.ParsePublicKey(decoded)
	if err != nil {
		return nil, fmt.Errorf("failed to parse.\n%v", err)
	}
	sshCryptoPublicKey, ok := sshPubKey.(ssh.CryptoPublicKey)
	if !ok {
		return nil, fmt.Errorf("invalid...")
	}

	return sshCryptoPublicKey.CryptoPublicKey(), nil
}

// func encryptMessage(plainMessage string, pubKey *rsa.PublicKey) (string, error) {
func encryptMessage(plainMessage string, pubKey crypto.PublicKey) (string, error) {
	switch publicKey := pubKey.(type) {
	case *rsa.PublicKey:
		// label := []byte("orderS")
		var label []byte

		// crypto/rand.Reader is a good source of entropy for randomizing the
		// encryption function.
		rng := rand.Reader
		ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, publicKey, []byte(plainMessage), label)
		if err != nil {
			return "", err
		}

		return string(ciphertext), nil
	case *ed25519.PublicKey:
		// Encrypt関数がない。https://pkg.go.dev/crypto/ed25519
		return "", fmt.Errorf("not supported")
	default:
		return "", fmt.Errorf("failed to encrypt")
	}
}
