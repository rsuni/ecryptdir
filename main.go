package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/atotto/clipboard"
)

type Config struct {
	//Secret key
	Key       string
	Directory string
}

func main() {

	//Commands
	createConfigFileCommand := flag.Bool("create-config", false, "Command Create config file create default config file.")
	encryptCommand := flag.Bool("e", false, "Encrypt files")
	decryptCommand := flag.Bool("d", false, "Decrypt files")
	generateKeyCommand := flag.Bool("k", false, "Generate key")

	//Settings
	configFileName := flag.String("config", "", "Config file include app settings.")

	flag.Parse()

	if *createConfigFileCommand {
		createConfigFile()
		return
	}
	if *generateKeyCommand {
		key, _ := GenerateRandomString(32)
		clipboard.WriteAll((key))
		log.Println("Key generated, show clipboard")
		return
	}

	if *configFileName == "" {
		fmt.Println("Error, config file must be entered")
		os.Exit(1)
	}

	if *encryptCommand {
		encrypt()
		return
	}
	if *decryptCommand {
		decrypt()
		return
	}
	return
}

func createConfigFile() {

	return
}
func encrypt() {
	log.Println("encrypted")
	return
}
func decrypt() {
	log.Println("decrypted")
	return
}

func EncryptText(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

//Decrypt symetric decrypt text by key
func DecryptText(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}

//GenerateRandomKey generate random key by definied length
func GenerateRandomBytes(length int) ([]byte, error) {

	key := make([]byte, length)
	_, err := rand.Read(key) //generate random key
	if err != nil {
		return key, err
	}

	return key, nil
}

func GenerateRandomString(length int) (string, error) {
	b, err := GenerateRandomBytes(length)
	key := base64.URLEncoding.EncodeToString(b)

	if len(key) > 32 {
		key = key[0:31]
	}
	return key, err
}
