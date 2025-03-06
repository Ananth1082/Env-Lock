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
	case "get":
		Get(args)
	case "help":
		Help()
	default:
		fmt.Println("Invalid command!")
		Help()
	}
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

func Get(args []string) {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	id := createCmd.Int64("id", 0, "ID of the file")

	// Parse the flags from args
	err := createCmd.Parse(args)
	if err != nil || *id == 0 {
		GetUsage()
	}

	file, err := DB.GetFile(*id)
	if err != nil {
		fmt.Println("File not found")
		return
	}
	fmt.Println("File Details:")
	fmt.Println("ID:", file.ID)
	fmt.Println("Name:", file.Name)
	fmt.Println("Description:", file.Description)

	decryptFile(file.ID)
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
