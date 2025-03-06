package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/argon2"
)

func deriveMasterKey(password string) ([]byte, []byte) {
	salt := make([]byte, ARGON2_SALT_SIZE)
	_, err := rand.Read(salt)
	if err != nil {
		fmt.Println("Error generating salt:", err)
	}
	return argon2.IDKey([]byte(password), salt, 1, ARGON2_MEM, 4, ARGON2_KEY_LEN), salt
}

func deriveMasterKeyWithSalt(password string, salt []byte) []byte {
	if len(salt) != ARGON2_SALT_SIZE {
		panic("Invalid salt size")
	}
	return argon2.IDKey([]byte(password), salt, 1, ARGON2_MEM, 4, ARGON2_KEY_LEN)
}

func generateKey() ([]byte, error) {
	key := make([]byte, AES_KEY_SIZE)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func encrypt(data []byte, key []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, NONCE_SIZE)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	ciphertext := aesGCM.Seal(nil, nonce, data, nil)
	return nonce, ciphertext, nil
}

func encryptTest() {
	fileData, err := os.ReadFile("test/test.env")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	aesKey, err := generateKey()
	if err != nil {
		fmt.Println("Error generating key:", err)
		return
	}

	nonce, encryptedData, err := encrypt(fileData, aesKey)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return
	}

	output, err := os.Create("test/test.enc")
	if err != nil {
		fmt.Println("Error creating encrypted file:", err)
		return
	}
	defer output.Close()

	output.Write(nonce)
	output.Write(encryptedData)

	fmt.Println("AES Key (hex):", hex.EncodeToString(aesKey))
	password := ""
	fmt.Println("Enter password:")
	fmt.Scanln(&password)
	masterKey, salt := deriveMasterKey(password)
	nonce, encryptedData, err = encrypt(aesKey, masterKey)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return
	}
	fmt.Println("Encrypted key: ", hex.EncodeToString(append(nonce, encryptedData...)))
	fmt.Println("Salt: ", hex.EncodeToString(salt))
}
