package main

import (
	"flag"
	"fmt"
	"os"
)

// GetCmd extracts the main command and subcommand from CLI args
func GetCmd() (string, []string) {
	if len(os.Args) < 2 {
		Help()
		os.Exit(1)
	}
	cmd := os.Args[1]
	args := os.Args[2:] // Pass remaining args (including flags)
	return cmd, args
}

// Control routes the command to the appropriate function
func Control(cmd string, args []string) {
	switch cmd {
	case "create":
		Create(args) // Pass args to create
	case "delete":
		Delete(args)
	case "list":
		List(args)
	case "help":
		Help()
	default:
		fmt.Println("Invalid command!")
		Help()
	}
}

// CreateUsage prints usage for the create command
func CreateUsage() {
	fmt.Println("Usage: elk-cli create -f [file_path] -n [name] -d [description]")
	os.Exit(1)
}

// Help prints general usage info
func Help() {
	fmt.Println("Usage: elk-cli [command] [flags]")
	fmt.Println("Commands:")
	fmt.Println("  create  - Create a new encrypted env file")
	fmt.Println("  delete  - Delete an env file")
	fmt.Println("  list    - List all env files")
	fmt.Println("  help    - Show usage information")
	os.Exit(1)
}

// Create handles the "create" command
func Create(args []string) {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	file := createCmd.String("f", "", "Path of env file")
	name := createCmd.String("n", "", "Name of the env file")
	description := createCmd.String("d", "Env file", "Description of the env file")

	// Parse the flags from args
	err := createCmd.Parse(args)
	if err != nil {
		CreateUsage()
	}

	if *file == "" {
		fmt.Println("Error: -f (file path) is required")
		CreateUsage()
	}

	// Simulate encryption function
	fileDetails := encryptFile(*file, *name, *description)
	fmt.Println("File Encrypted:")
	fmt.Println("ID:", fileDetails.ID)
	fmt.Println("Name:", fileDetails.Name)
	fmt.Println("Description:", fileDetails.Description)
}

// Delete command
func Delete(args []string) {
	fmt.Println("Delete functionality not implemented yet")
}

// List command
func List(args []string) {
	fmt.Println("List functionality not implemented yet")
}

// CLI function to handle command input
func CLI() {
	cmd, args := GetCmd()
	Control(cmd, args)
}
