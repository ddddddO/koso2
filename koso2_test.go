package koso2

import (
	"fmt"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	tests := map[string]struct {
		ghUserID     string
		plainMessage string
		callbacks    []callbackFn
	}{
		"one callback": {
			ghUserID:     "ddddddO",
			plainMessage: "aaaaAAaあaa",
			callbacks: []callbackFn{
				func(encrypted string) error {
					fmt.Print(encrypted)
					return nil
				},
			},
		},
		"multiple callbacks": {
			ghUserID:     "ddddddO",
			plainMessage: "aaaaAAaあaa",
			callbacks: []callbackFn{
				func(encrypted string) error {
					fmt.Println("callback 1")
					fmt.Print(encrypted)
					return nil
				},
				func(encrypted string) error {
					fmt.Println("callback 2")
					fmt.Print(encrypted)
					return nil
				},
			},
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if err := Run(tt.ghUserID, tt.plainMessage, tt.callbacks...); err != nil {
				t.Errorf("failed to test:\n%+v\n", err)
			}
		})
	}
}

func TestRunConcurrently(t *testing.T) {
	tests := map[string]struct {
		ghUserID     string
		plainMessage string
		callbacks    []callbackFn
	}{
		"one callback": {
			ghUserID:     "ddddddO",
			plainMessage: "aaaaAAaあaa",
			callbacks: []callbackFn{
				func(encrypted string) error {
					fmt.Print(encrypted)
					return nil
				},
			},
		},
		"multiple callbacks": {
			ghUserID:     "ddddddO",
			plainMessage: "aaaaAAaあaa",
			callbacks: []callbackFn{
				func(encrypted string) error {
					time.Sleep(time.Second)
					fmt.Println("callback 1")
					fmt.Print(encrypted)
					return nil
				},
				func(encrypted string) error {
					time.Sleep(time.Second)
					fmt.Println("callback 2")
					fmt.Print(encrypted)
					return nil
				},
				func(encrypted string) error {
					time.Sleep(time.Second)
					fmt.Println("callback 3")
					fmt.Print(encrypted)
					return nil
				},
				func(encrypted string) error {
					time.Sleep(time.Second)
					fmt.Println("callback 4")
					fmt.Print(encrypted)
					return nil
				},
			},
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if err := RunConcurrently(tt.ghUserID, tt.plainMessage, tt.callbacks...); err != nil {
				t.Errorf("failed to test:\n%+v\n", err)
			}
		})
	}
}

func TestParsePublicKeyAndEncryptMessage(t *testing.T) {
	tests := map[string]struct {
		pk string
	}{
		"rsa": {pk: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDFWREDexIn+DjpQRVxpWrOLrl6z+bNFl6iohm55FaEkiPyuVHT9FXYTwPetVBLcRRF7GCDhwllSW7IT5Ucn+EsWyJeLUUEk7nGZHSdZ/4ssUTekl9aSxZz3aCD702mlu3hj2ohcnFkkYQQ507vh8TwfHstO7tbC5iXO9eHJCmzgcuX0pRNzljqXkrv/k97smFUu3uLupeYiDNJTMz9pAxSaDlUZG/T5lUDa0qizcCfGIayJjPy2SruwKKqP7lLdK4JwFeCT/ibqwdWEL//Wg4C19imyqdZpcHrn8vexgaYKpWrjwbFDgLB9xtwKhsTlXltBEuMPGk5Cqz10g2Bgso7"},
		// "ed25519": {pk: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGtjunUgSSpmgurFspc5gMtkTvK5owMz0I9TZ6zeNavR"},
	}

	plainMessage := "aaaCCCC"

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			gotPubKey, gotErr := parsePublicKey(tt.pk)
			if gotErr != nil {
				t.Errorf("failed to test:\n%+v\n", gotErr)
			}
			if _, gotErr := encryptMessage(plainMessage, gotPubKey); gotErr != nil {
				t.Errorf("failed to test:\n%+v\n", gotErr)
			}
		})
	}
}
