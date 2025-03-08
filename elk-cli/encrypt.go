package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
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

func encryptFile(fileName string, outFileName string) ([]byte, []byte) {

	fileData, err := os.ReadFile(fileName)
	if err != nil {
		log.Panicln("Error reading file:", err)
		return nil, nil
	}

	aesKey, err := generateKey()
	if err != nil {
		log.Panicln("Error generating key:", err)
		return nil, nil
	}

	nonce, encryptedData, err := encrypt(fileData, aesKey)
	if err != nil {
		log.Panicln("Encryption error:", err)
		return nil, nil
	}
	output, err := os.Create(outFileName)
	if err != nil {
		log.Panicln("Error creating encrypted file:", err)
		return nil, nil
	}
	defer output.Close()
	output.Write(nonce)
	output.Write(encryptedData)
	fmt.Print("Password: ")
	password := EnterPassword()
	check := CheckPassword(password)
	if !check {
		log.Fatalln("Incorrect password")
	}
	masterKey, salt := deriveMasterKey(string(password))
	nonce, encaesKey, err := encrypt(aesKey, masterKey)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return nil, nil
	}

	return append(nonce, encaesKey...), salt
}
