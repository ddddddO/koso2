package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

// NOTE: 一旦、RSAのみ考える
func main() {
	const ghUserID = "ddddddO"
	const plain = `AAA
こんにちは
！！!!
`

	cb := func(enc string) error {
		fmt.Print(enc)
		return nil
	}

	if err := run(ghUserID, plain, cb); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ghUserID, plainMessage string, callback func(encrypted string) error) error {
	pubKeys, err := fetchPublicKeys(ghUserID)
	if err != nil {
		return err
	}

	pk := pubKeys[0] // ローカルにあるのがこれなので、一旦
	rsaPubKey, err := parsePublicKey(pk)
	if err != nil {
		return err
	}

	// TODO: 続きは暗号化から
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
		Host: "github.com",
		Path: fmt.Sprintf("%s.keys", ghUserID),
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
		rawPub := scanner.Text()		
		pubKeys = append(pubKeys, strings.TrimLeft(rawPub, "ssh-rsa "))
	}

	return pubKeys, nil
}

func parsePublicKey(publicKey string) (*rsa.PublicKey, error) {
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

	pub := sshCryptoPublicKey.CryptoPublicKey()
	return pub.(*rsa.PublicKey), nil
}

func encryptMessage(plainMessage string, pubKey *rsa.PublicKey) (string, error) {
	// label := []byte("orderS")
	var label []byte 

	// crypto/rand.Reader is a good source of entropy for randomizing the
	// encryption function.
	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, pubKey, []byte(plainMessage), label)
	if err != nil {
		return "", fmt.Errorf("Error from encryption: %s\n", err)
	}

	return string(ciphertext), nil
}