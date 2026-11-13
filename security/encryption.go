package security

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	gochoice "github.com/TwiN/go-choice"
	"github.com/secure-systems-lab/go-securesystemslib/encrypted"
)

const (
	FAILED_DECRYPTION_MESSAGE string = "Failed to decrypt"
)

var Encrypt = make(chan string)
var Decrypt = make(chan []byte)

// Creates the file and store passphrase within it
func createPassphraseFile(passphrasePath string, passphrase string) {
	folder, _, _ := strings.Cut(passphrasePath, "passphrase.txt")
	if err := os.MkdirAll(folder, 0700); err != nil {
		fmt.Printf("Error during creation of passphrase file. Error: %v\n", err)
	} else {
		file, err := os.Create(passphrasePath)
		if err != nil {
			fmt.Printf("Error during creation of passphrase file. Error: %v\n", err)
		} else {
			if _, err := file.Write([]byte(passphrase)); err != nil {
				fmt.Printf("Error during storing passphrase. Error: %v\n", err)
			}
		}
	}
}

// Grab the saved passphrase
func grabPassphraseSaved(passphrasePath string) string {
	bytes, err := os.ReadFile(passphrasePath)
	if err != nil {
		fmt.Printf("Couldn't read the passphrase saved. Error: %v\n", err)
	} else {
		return string(bytes)
	}
	return ""
}

// Sets up the encryption by applying passphrasing
func Setup() {
	passphraseLocation, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("Error during finding passphrase save location: %v\n", err)
	}
	passphraseLocation = path.Join(passphraseLocation, "clipboard-sync", "passphrase.txt")

	var passphrase string
	_, err = os.Stat(passphraseLocation)
	if errors.Is(err, os.ErrNotExist) {
		passphrase = askForPassphrase()
		savePassphrase := askShouldSavePassphrase()
		if savePassphrase == "Yes" {
			createPassphraseFile(passphraseLocation, passphrase)
		}
	} else {
		passphrase = grabPassphraseSaved(passphraseLocation)
	}

	passphraseBytes := []byte(passphrase)

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

// Prompts the user to see if they want to save the passphrase
func askShouldSavePassphrase() string {
	choice, _, err := gochoice.Pick(
		"Do you want to save your passphrase?",
		[]string{
			"Yes",
			"No",
		},
		gochoice.OptionBackgroundColor(gochoice.Black),
		gochoice.OptionTextColor(gochoice.White),
		gochoice.OptionSelectedTextColor(gochoice.Red),
		gochoice.OptionSelectedTextBold(),
	)
	if err != nil {
		return "No"
	}

	return choice
}
