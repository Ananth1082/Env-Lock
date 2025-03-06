package main

import (
	"fmt"
	"os"
)

func CreateUsage() {
	fmt.Println("Usage: elk-cli create -f [file_path] -n [name] -d [description]")
	os.Exit(1)
}

func GetUsage() {
	fmt.Println("Usage: elk-cli get -id [file_id]")
	os.Exit(1)
}

func Help() {
	fmt.Println("Usage: elk-cli [command] [flags]")
	fmt.Println("Commands:")
	fmt.Println("  create  - Create a new encrypted env file")
	fmt.Println("  delete  - Delete an env file")
	fmt.Println("  list    - List all env files")
	fmt.Println("  help    - Show usage information")
	os.Exit(1)
}
