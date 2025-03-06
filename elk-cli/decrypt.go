package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path"
)

func decrypt(encryptedData []byte, nonce []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesGCM.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func decryptTest() {
	encFile, err := os.ReadFile("test/test.enc")
	if err != nil {
		fmt.Println("Error reading encrypted file:", err)
		return
	}

	var aesKeyHex string
	var password string
	var salt string

	fmt.Print("Enter AES encrypted Key (hex): ")
	fmt.Scanln(&aesKeyHex)

	fmt.Println("Enter the password: ")
	fmt.Scanln(&password)

	fmt.Println("Enter the salt: ")
	fmt.Scanln(&salt)

	aesEncKey, err := hex.DecodeString(aesKeyHex)
	if err != nil {
		fmt.Println("Invalid AES key")
		return
	}

	saltBytes, err := hex.DecodeString(salt)
	if err != nil || len(saltBytes) != ARGON2_SALT_SIZE {
		fmt.Println("Invalid salt")
		return
	}
	masterKey := deriveMasterKeyWithSalt(password, saltBytes)
	aesKey, err := decrypt(aesEncKey[NONCE_SIZE:], aesEncKey[:NONCE_SIZE], masterKey)
	if err != nil {
		log.Panicf("Error: %v", err)
	}

	if len(encFile) < NONCE_SIZE {
		fmt.Println("Invalid encrypted file")
		return
	}

	nonce := encFile[:NONCE_SIZE]
	encryptedData := encFile[NONCE_SIZE:]

	decryptedData, err := decrypt(encryptedData, nonce, aesKey)
	if err != nil {
		fmt.Println("Decryption error:", err)
		return
	}

	err = os.WriteFile(path.Join(ENC_DIR, "test.enc"), decryptedData, 0644)
	if err != nil {
		fmt.Println("Error writing decrypted file:", err)
		return
	}

	fmt.Println("Decryption successful! File saved as test_decrypted.env")
}
