package main

import (
	"bufio"
	"elk/elk/util"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	_ "embed"

	"github.com/pelletier/go-toml/v2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

const (
	PSWD_HASH_PRICE = 12
)

//go:embed user.toml
var init_temp string

func HashPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, PSWD_HASH_PRICE)
}

func ElkInit() {
	_, err := os.Stat(CONFIG_FILE)
	if err == nil {
		util.PrintWarning("ELK CLI already initialized")
		return
	}
	util.PrintSuccess("Initializing ELK CLI...")
	util.PrintWarning(`
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
		util.PrintError("Error creating config directory:")
		return
	}
	util.PrintSuccess("1) Created config directory!!!")

	if err := os.Mkdir(ENC_DIR, 0700); err != nil && !os.IsExist(err) {
		util.PrintError("Error creating encryption directory")
		return
	}
	util.PrintSuccess("2) Created encryption directory!!!")
	reader := bufio.NewReader(os.Stdin)

	util.PrintPrompt("Enter your name: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		util.PrintError("Error reading name")
		return
	}
	username = strings.TrimSpace(username)

	util.PrintPrompt("Enter your email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		util.PrintError("Error reading email")
		return
	}
	email = strings.TrimSpace(email)

	util.PrintPrompt("3) Enter your password (remember this password, you will need it to access the CLI):")
	bytePassword := EnterPassword()
	bcryptHash, err := HashPassword(bytePassword)
	if err != nil {
		util.PrintError("Error hashing password")
		return
	}

	file, err := os.Create(CONFIG_FILE)
	if err != nil {
		util.PrintError("Error creating config file")
		return
	}
	defer file.Close()

	configData := fmt.Sprintf(init_temp, username, email, hex.EncodeToString(bcryptHash))
	if _, err := file.WriteString(configData); err != nil {
		util.PrintError("Error writing to config file")
		return
	}

	if err := file.Chmod(0600); err != nil {
		fmt.Println("Error setting file permissions:", err)
		return
	}
	InitDB()
	util.PrintSuccess("Created config file!!!")
}

func EnterPassword() []byte {
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatalln("\nError reading password:", err)
	}
	fmt.Println()
	return bytePassword
}

func CheckPassword(password []byte) bool {
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
	pswdHash, err := hex.DecodeString(config.Owner.Password_hash)
	if err != nil {
		log.Fatal("Corrupted password hash")
	}
	err = bcrypt.CompareHashAndPassword(pswdHash, password)
	if err != nil {
		util.PrintError("Incorrect password")
		return false
	}
	fmt.Println("Correct password")
	return true
}

func CheckUser() bool {
	data, err := os.ReadFile(CONFIG_FILE)
	if err != nil {
		util.PrintError("ELK CLI not initialized")
		ElkInit()
		return false
	}
	var User struct {
		Owner struct {
			Name          string
			Email         string
			Password_hash string
		}
	}
	err = toml.Unmarshal(data, &User)
	if err != nil || User.Owner.Name == "" || User.Owner.Email == "" || User.Owner.Password_hash == "" {
		log.Println("Invalid config file", err)
		util.PrintError("Invalid config file")
		ElkInit()
		return false
	}
	return true
}
