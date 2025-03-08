package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
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

func decryptFile(fileId int64, outFile string) error {
	file, err := DB.GetFile(fileId)
	if err != nil {
		return ErrFileNotFound
	}
	encFile, err := os.ReadFile(path.Join(ENC_DIR, file.Path))
	if err != nil {
		return ErrFileNotFound
	}

	fmt.Print("Password: ")
	pswd := EnterPassword()
	CheckPassword(pswd)
	aesEncKey, err := hex.DecodeString(file.Key)
	if err != nil {
		fmt.Println("Invalid AES key")
		return ErrInvalidKeyFormat
	}

	saltBytes, err := hex.DecodeString(file.Salt)
	if err != nil || len(saltBytes) != ARGON2_SALT_SIZE {
		fmt.Println("Invalid salt")
		return ErrInvalidKeyFormat
	}
	masterKey := deriveMasterKeyWithSalt(string(pswd), saltBytes)
	aesKey, err := decrypt(aesEncKey[NONCE_SIZE:], aesEncKey[:NONCE_SIZE], masterKey)
	if err != nil {
		return ErrDecryption
	}

	if len(encFile) < NONCE_SIZE {
		fmt.Println("Invalid encrypted file")
		return ErrInvalidEncryptionFile
	}

	nonce := encFile[:NONCE_SIZE]
	encryptedData := encFile[NONCE_SIZE:]

	decryptedData, err := decrypt(encryptedData, nonce, aesKey)
	if err != nil {
		fmt.Println("Decryption error:", err)
		return ErrDecryption
	}

	if outFile == "" {
		outFile = ".env"
	}
	err = os.WriteFile(outFile, decryptedData, 0644)
	if err != nil {
		return ErrIo
	}
	return nil
}
