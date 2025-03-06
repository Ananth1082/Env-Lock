package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/pelletier/go-toml/v2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

const (
	PSWD_HASH_PRICE = 12
)

const init_temp = ` # ELK CLI configuration file
# This file contains the configuration for the ELK CLI tool
                                                                             
title = "elk config"

[owner]
name = "%s"
email = "%s"
password_hash = "%s"
`

func HashPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, PSWD_HASH_PRICE)
}

func ElkInit() {
	_, err := os.Stat(CONFIG_FILE)
	if err == nil {
		fmt.Println("ELK CLI already initialized")
		return
	}

	fmt.Print(`Initializing ELK CLI...
	
 /$$$$$$$$                            /$$                           /$$      
| $$_____/                           | $$                          | $$      
| $$       /$$$$$$$  /$$    /$$      | $$        /$$$$$$   /$$$$$$$| $$   /$$
| $$$$$   | $$__  $$|  $$  /$$/      | $$       /$$__  $$ /$$_____/| $$  /$$/
| $$__/   | $$  \ $$ \  $$/$$/       | $$      | $$  \ $$| $$      | $$$$$$/ 
| $$      | $$  | $$  \  $$$/        | $$      | $$  | $$| $$      | $$_  $$ 
| $$$$$$$$| $$  | $$   \  $/         | $$$$$$$$|  $$$$$$/|  $$$$$$$| $$ \  $$
|________/|__/  |__/    \_/          |________/ \______/  \_______/|__/  \__/

`)

	if err := os.Mkdir(CONFIG_DIR, 0700); err != nil && !os.IsExist(err) {
		fmt.Println("Error creating config directory:", err)
		return
	}
	fmt.Println("Created config directory!!!")

	if err := os.Mkdir(ENC_DIR, 0700); err != nil && !os.IsExist(err) {
		fmt.Println("Error creating encryption directory:", err)
		return
	}
	fmt.Println("Created encryption directory!!!")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your name: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading name:", err)
		return
	}
	username = strings.TrimSpace(username)

	fmt.Print("Enter your email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading email:", err)
		return
	}
	email = strings.TrimSpace(email)

	fmt.Print("Enter your password (Remember the password): ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("\nError reading password:", err)
		return
	}
	fmt.Println()

	bcryptHash, err := HashPassword(bytePassword)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return
	}

	file, err := os.Create(CONFIG_FILE)
	if err != nil {
		fmt.Println("Error creating config file:", err)
		return
	}
	defer file.Close()

	configData := fmt.Sprintf(init_temp, username, email, hex.EncodeToString(bcryptHash))
	if _, err := file.WriteString(configData); err != nil {
		fmt.Println("Error writing to config file:", err)
		return
	}

	if err := file.Chmod(0600); err != nil {
		fmt.Println("Error setting file permissions:", err)
		return
	}

	fmt.Println("ELK CLI initialized successfully!")
}

func CheckPassword(password string) bool {
	data, err := os.ReadFile(CONFIG_FILE)
	if err != nil {
		log.Fatal(err)
	}
	config := new(struct {
		Owner struct {
			Password_hash string
			Name          string
		}
	})
	err = toml.Unmarshal(data, config)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Config : ", config)
	pswdHash, err := hex.DecodeString(config.Owner.Password_hash)
	if err != nil {
		log.Fatal("Corrupted password hash")
	}
	err = bcrypt.CompareHashAndPassword(pswdHash, []byte(password))
	if err != nil {
		fmt.Println("Incorrect password")
		return false
	}
	fmt.Println("Correct password")
	return true
}
