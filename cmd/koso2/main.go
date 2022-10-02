package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/ddddddO/koso2"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "Specify only GitHub userID.")
		os.Exit(1)
	}
	ghUserID := args[1]

	dst := &bytes.Buffer{}
	io.Copy(dst, os.Stdin)
	plainMessage := dst.String()

	outputEncryptedMessage := func(encrypted string) error {
		fmt.Print(encrypted)
		return nil
	}

	if err := koso2.Run(ghUserID, plainMessage, outputEncryptedMessage); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
