package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/google/uuid"
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

func encryptFile(file string) *File {
	fileData, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	aesKey, err := generateKey()
	if err != nil {
		fmt.Println("Error generating key:", err)
		return nil
	}

	nonce, encryptedData, err := encrypt(fileData, aesKey)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return nil
	}

	encFile := uuid.New().String() + ".enc"
	output, err := os.Create(path.Join(ENC_DIR, encFile))
	if err != nil {
		fmt.Println("Error creating encrypted file:", err)
		return nil
	}
	defer output.Close()

	output.Write(nonce)
	output.Write(encryptedData)

	password := ""
	fmt.Println("Enter password:")
	fmt.Scanln(&password)
	check := CheckPassword(password)
	if !check {
		log.Fatalln("Incorrect password")
	}
	masterKey, salt := deriveMasterKey(password)
	nonce, encryptedData, err = encrypt(aesKey, masterKey)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return nil
	}
	fileMetaData := File{
		Name:        file,
		Description: "Encrypted file",
		Path:        encFile,
		Key:         hex.EncodeToString(append(nonce, encryptedData...)),
		Salt:        hex.EncodeToString(salt),
	}
	newfile, err := DB.CreateFile(&fileMetaData)
	if err != nil {
		log.Fatal(err)
	}
	return newfile
}
