package main

import "errors"

var (
	ErrFileNotFound          = errors.New("file not found")
	ErrInvalidEncryptionFile = errors.New("invalid encryption")
	ErrInvalidKeyFormat      = errors.New("invalid key format")
	ErrInvalidPassword       = errors.New("invalid password")
	ErrDecryption            = errors.New("error decrypting file")
	ErrIo                    = errors.New("error doing io operation")
	ErrGeneratingKey         = errors.New("error generating key")
	ErrEncryption            = errors.New("error encrypting file")
)
