package security

import (
	"fmt"

	"github.com/secure-systems-lab/go-securesystemslib/encrypted"
)

const (
	FAILED_DECRYPTION_MESSAGE string = "Failed to decrypt"
)

var Encrypt = make(chan string)
var Decrypt = make(chan []byte)

// Sets up the encryption by applying passphrasing
func Setup() {
	passphrase := askForPassphrase()
	passphraseBytes := []byte(passphrase)
	// fmt.Printf("Verified your passphrase as: %s\n", passphrase)

	go func() {
		for text := range Encrypt {
			bytes := []byte(text)
			encryptedBytes, _ := encryptBytes(bytes, passphraseBytes)

			Encrypt <- string(encryptedBytes)
		}
	}()

	go func() {
		for bytes := range Decrypt {
			decryptedBytes, err := decryptBytes(bytes, passphraseBytes)
			if err != nil {
				Decrypt <- []byte(FAILED_DECRYPTION_MESSAGE)
			} else {
				Decrypt <- decryptedBytes
			}
		}
	}()
}

// Encrypts the bytes given to it using a passphrase
func encryptBytes(bytes []byte, passphrase []byte) ([]byte, error) {
	encryptedBytes, err := encrypted.Encrypt(bytes, passphrase)
	if err != nil {
		fmt.Println("Error during encryption\n", err)
	}

	return encryptedBytes, err
}

// Decrypts the ciphertext given to it by a using a passphrase
func decryptBytes(ciphertext []byte, passphrase []byte) ([]byte, error) {
	decryptedBytes, err := encrypted.Decrypt(ciphertext, passphrase)
	if err != nil {
		fmt.Println("Error during decryption\n", err)
	}

	return decryptedBytes, err
}

// Prompts the user to enter a passphrase
func askForPassphrase() string {
	var passphrase string

	fmt.Print("Type your passphrase: ")
	fmt.Scan(&passphrase)
	if len(passphrase) >= 8 {
		return passphrase
	} else {
		fmt.Println("Your passphrase must be atleast 8 characters.")
		return askForPassphrase()
	}
}
