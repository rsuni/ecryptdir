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
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
)

const VERSION = "1.0"

type Config struct {
	//Secret key
	Key       string
	Directory string
}

var config Config

func main() {

	//Commands
	ccCommand := flag.Bool("cc", false, "Create config command: create default config file.")
	eCommand := flag.Bool("e", false, "Encrypt files command")
	dCommand := flag.Bool("d", false, "Decrypt files command")
	gkCommand := flag.Bool("gk", false, "Generate key command")
	vCommand := flag.Bool("v", false, "App version")

	//Settings
	configFileName := flag.String("config", "./config.json", "Config file name include app settings.")
	fileFindRule := flag.String("find", "", "What find f.e. file1*")

	flag.Parse()

	if *vCommand {
		fmt.Printf("Version %s \n", VERSION)
		return
	}

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

	config = getConfig()

	if config.Directory == "" {
		fmt.Println("Error, directory is not specified")
		os.Exit(1)
	}
	if config.Key == "" {
		fmt.Println("Error, key is not specified")
		os.Exit(1)
	}

	if *fileFindRule == "" {
		fmt.Println("Error, what find is not specified")
		os.Exit(1)
	}
	files, _ := FindFiles(config.Directory, *fileFindRule)
	_ = files

	if *eCommand {
		encrypt(files)
		return
	}
	if *dCommand {
		decrypt(files)
		return
	}
	return
}

func FindFiles(directory, whatFind string) ([]string, error) {
	filesFinded := []string{}

	last := string(whatFind[len(whatFind)-1])
	if last == "*" {
		whatFind = whatFind[:len(whatFind)-1]
	}

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if !f.IsDir() && strings.HasPrefix(f.Name(), whatFind) {
			filesFinded = append(filesFinded, filepath.Join(directory, f.Name()))
		}
	}

	return filesFinded, nil
}

func getConfig() Config {

	raw, err := readFile("./config.json")
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
func encrypt(files []string) {

	for _, file := range files {
		encryptSingleFile(file)
	}

	log.Println("encrypted")
	return
}
func encryptSingleFile(fileName string) {
	raw, err := readFile(fileName)
	checkError(err)

	result, err := encryptText([]byte(config.Key), raw)
	checkError(err)

	err = ioutil.WriteFile(fileName, result, 0644)
	checkError(err)

	return
}

func decrypt(files []string) {
	for _, file := range files {
		decryptSingleFile(file)
	}

	log.Println("decrypted")
	return
}
func decryptSingleFile(fileName string) {
	raw, err := readFile(fileName)
	checkError(err)

	result, err := decryptText([]byte(config.Key), raw)
	checkError(err)

	err = ioutil.WriteFile(fileName, result, 0644)
	checkError(err)

	return
}
func encryptText(key, text []byte) ([]byte, error) {
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
func decryptText(key, text []byte) ([]byte, error) {
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

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return
}
func readFile(fileName string) ([]byte, error) {
	raw, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, nil
	}
	return raw, nil
}
