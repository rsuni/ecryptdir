package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
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
	ccCommand := flag.Bool("cc", false, "Create config command: create default config file.")
	eCommand := flag.Bool("e", false, "Encrypt files command")
	dCommand := flag.Bool("d", false, "Decrypt files command")
	gkCommand := flag.Bool("gk", false, "Generate key command")

	//Settings
	configFileName := flag.String("config", "./config.json", "Config file include app settings.")
	//key := flag.String("key", "", "Key")

	flag.Parse()

	if *ccCommand {
		createConfigFile()
		return
	}
	if *gkCommand {
		key, _ := GenerateRandomString(32)
		clipboard.WriteAll((key))
		log.Println("Key generated, show clipboard")
		return
	}

	if *configFileName == "" {
		fmt.Println("Error, config file must be entered")
		os.Exit(1)
	}

	if _, err := os.Stat(*configFileName); os.IsNotExist(err) {
		fmt.Println("Error, config file does not exists")
		os.Exit(1)
	}

	config := getConfig()

	if config.Directory == "" {
		fmt.Println("Error, directory is not specified")
		os.Exit(1)
	}

	if *eCommand {
		encrypt()
		return
	}
	if *dCommand {
		decrypt()
		return
	}
	return
}

func getConfig() Config {
	raw, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c Config
	json.Unmarshal(raw, &c)
	return c
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

	if len(key) > length {
		key = key[0:length]
	}
	return key, err
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return
}
