package main

import (
	"os"
	"path"
)

const (
	ARGON2_SALT_SIZE = 16
	ARGON2_MEM       = 64 * 1024
	ARGON2_KEY_LEN   = 32
	NONCE_SIZE       = 12
	AES_KEY_SIZE     = 32
)

var (
	HOME_DIR, _ = os.UserHomeDir()
	CONFIG_DIR  = path.Join(HOME_DIR, ".elk")
	ENC_DIR     = path.Join(CONFIG_DIR, "enc")
	CONFIG_FILE = path.Join(CONFIG_DIR, "config.toml")
)
