# Koso x Koso

Encrypt message using public key of user with specified GitHub id.</br>
Encrypted message can be processed by your Go program.</br>
Inspired by **[naisho](https://github.com/moznion/naisho)** !

⚠Currently, only RSA is supported.

## Usage

### Encrypt

```go
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/ddddddO/koso2"
)

func main() {
	ghUserID := "ddddddO"

	plainMessage := `こんにちは！
Hello!
Halo!
Bonjour!
你好!
안녕하세요!
habari!
`

    // 'encrypted' is the encrypted plainMessage.
	callback := func(encrypted string) error {
    //
    // For example, we could add processing to send encrypted message to Slack.
    //
    fmt.Print(encrypted)
		return nil
	}

	if err := koso2.Run(ghUserID, plainMessage, callback); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```

### Decrypt

```sh
$ openssl pkeyutl -decrypt -in <Encrypted Message File> -inkey ~/.ssh/id_rsa \
    -pkeyopt rsa_padding_mode:oaep -pkeyopt rsa_oaep_md:sha256
```

## Install

### Package
```sh
$ go get github.com/ddddddO/koso2
```

### CLI

Currently, this CLI only has function of encrypting message with public key of specified GitHub ID and outputting it to standard output.

```sh
$ go install github.com/ddddddO/koso2/cmd/koso2@latest
```
