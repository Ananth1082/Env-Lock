package main

import (
	"fmt"
	"os"
)

func CreateUsageExit() {
	fmt.Println("Usage: elk create -f [file_path] -n [name] -d [description]")
	os.Exit(1)
}

func GetUsageExit() {
	fmt.Println("Usage: elk get -id [file_id]")
	os.Exit(1)
}

func DeleteUsageExit() {
	fmt.Println("Usage: elk delete -id [file_id]")
	os.Exit(1)
}

func HelpExit() {
	fmt.Println("Usage: elk [command] [flags]")
	fmt.Println("Commands:")
	fmt.Println("  create  - Create a new encrypted env file")
	fmt.Println("  delete  - Delete an env file")
	fmt.Println("  list    - List all env files")
	fmt.Println("  help    - Show usage information")
	os.Exit(1)
}
