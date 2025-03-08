package main

import (
	"fmt"
)

func CreateUsage() {
	fmt.Println("Usage: elk create -f [file_path] -n [name] -d [description]")
}

func GetUsage() {
	fmt.Println("Usage: elk get -id [file_id]")
}

func DeleteUsage() {
	fmt.Println("Usage: elk delete -id [file_id]")
}

func Help() {
	fmt.Println("Usage: elk [command] [flags]")
	fmt.Println("Commands:")
	fmt.Println("  create  - Create a new encrypted env file")
	fmt.Println("  delete  - Delete an env file")
	fmt.Println("  list    - List all env files")
	fmt.Println("  help    - Show usage information")
}
