package main

import (
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key, err := GenerateRandomString(32) //16,24,32
	CheckErrorTest(t, err)
	text, err := GenerateRandomString(32)
	CheckErrorTest(t, err)
	key2 := []byte(key)
	text2 := []byte(text)

	result1, err := EncryptText(key2, text2)
	CheckErrorTest(t, err)

	result2, err := DecryptText(key2, result1)
	CheckErrorTest(t, err)

	text3 := string(result2)

	if text != text3 {

		t.Fatal("Encrypt decrypt test fail")
	}

	t.Log("Encrypt decrypt test ok")
	return
}

func TestGetConfig(t *testing.T) {
	config := getConfig()

	if config.Key == "" {
		t.Fatal("config doesnt contain key value.")
	}

	return
}
func TestFindFiles(t *testing.T) {
	config := getConfig()

	files, err := FindFiles(config.Directory, "file*")
	CheckErrorTest(t, err)

	if len(files) == 0 {
		t.Fatal("dont find files.")
	}

	return
}

func CheckErrorTest(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
	return
}
