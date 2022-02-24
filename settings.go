package main

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

type AccountAuthSecrets struct {
	Private *rsa.PrivateKey
	Public  *rsa.PublicKey
}

type SecretRepository struct {
	AccountAuth AccountAuthSecrets
}

type Settings struct {
	PostgresURL string `json:"postgresUrl"`
	Secrets     SecretRepository
}

func handleSecret[T any](private bool, name string) *T {
	bytes, readError := ioutil.ReadFile("secrets/" + name)
	if readError != nil {
		log.Printf("Failed to read '%v' in dir 'secrets.' Is it present? Exiting...\n", name)
		return nil
	}

	var key any
	if private {
		pKey, keyErr := ssh.ParseRawPrivateKey(bytes)

		if keyErr != nil {
			log.Printf("Invalid RSA private key in '%v.' Please correct it and start the application again. Exiting...", name)
			return nil
		}

		key = pKey
	} else {
		pKey, _, _, _, err1 := ssh.ParseAuthorizedKey(bytes)
		if err1 != nil {
			log.Printf("Could not parse RSA public key in '%v.' Please correct it and start the application again. Exiting...", name)
			return nil
		}

		sshKey, err2 := ssh.ParsePublicKey(pKey.Marshal())
		if err2 != nil {
			log.Printf("Could not parse RSA public key in '%v.' Please correct it and start the application again. Exiting...", name)
			return nil
		}

		cryptoKey := sshKey.(ssh.CryptoPublicKey)
		key = cryptoKey.CryptoPublicKey().(*rsa.PublicKey)
	}
	pass, _ := key.(*T)
	return pass
}

func ReadSettings() *Settings {
	// bytes read from the settings file
	var settingsContent []byte

	// create file if it doesn't exist
	if _, settingsError := os.Stat("settings.json"); errors.Is(settingsError, os.ErrNotExist) {
		created, createdError := os.Create("settings.json")
		if createdError != nil {
			log.Panic("An error occurred during creation of 'settings.json' file. Try creating it manually.")
			return nil
		}
		settingsContent = []byte(`{
  "postgresUrl": "postgres://user:pass@address:port/<database_name>"
}`)
		created.Write(settingsContent)
		created.Close()
		log.Println("No 'settings' file found. One was created for you, please modify it accordingly. Exiting...")
		time.Sleep(time.Second * 3)
		return nil
	}

	// fetch settings from file
	contentRead, _ := ioutil.ReadFile("settings.json")
	settingsContent = contentRead

	// unmarshal json into our var with corresponding struct
	var settings Settings
	json.Unmarshal(settingsContent, &settings)

	// generate secrets accordingly
	os.Mkdir("secrets", 0644)

	privateAuthKey := handleSecret[rsa.PrivateKey](true, "account_auth")
	publicAuthKey := handleSecret[rsa.PublicKey](false, "account_auth.pub")

	if privateAuthKey == nil || publicAuthKey == nil {
		time.Sleep(time.Second * 3)
		return nil
	}

	secrets := &settings.Secrets
	secrets.AccountAuth = AccountAuthSecrets{Private: privateAuthKey, Public: publicAuthKey}

	// return the settings accordingly
	return &settings
}
