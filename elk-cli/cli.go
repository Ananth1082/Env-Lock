package main

import (
	"elk/elk/util"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/google/uuid"
)

func GetCmd() (string, []string) {
	if len(os.Args) < 2 {
		Help()
		os.Exit(1)
	}
	cmd := os.Args[1]
	args := os.Args[2:]
	return cmd, args
}

func Control(cmd string, args []string) {
	switch cmd {
	case "create":
		Create(args)
	case "delete":
		Delete(args)
	case "list":
		List(args)
	case "get":
		Get(args)
	case "update":
		Update(args)
	case "help":
		Help()
	default:
		fmt.Println("Invalid command!")
		Help()
	}
}

func Create(args []string) {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	file := createCmd.String("f", "", "Path of env file")
	name := createCmd.String("n", "", "Name of the env file")
	description := createCmd.String("d", "Env file", "Description of the env file")

	err := createCmd.Parse(args)
	if err != nil {
		CreateUsage()
	}

	if *file == "" {
		fmt.Println("Error: -f (file path) is required")
		CreateUsage()
	}
	pathId := uuid.New().String() + ".enc"
	encFilepath := path.Join(ENC_DIR, pathId)
	key, salt := encryptFile(*file, encFilepath)

	fileDetails := &File{}
	fileDetails.Details.Name = *name
	fileDetails.Path = pathId
	fileDetails.Salt = hex.EncodeToString(salt)
	fileDetails.Key = hex.EncodeToString(key)
	fileDetails.Details.Description = *description
	err = DB.CreateFile(fileDetails)
	if err != nil {
		fmt.Println("Error: creating file in local db", err)
	}
	fmt.Println("File Encrypted:")
	fmt.Println("ID:", fileDetails.Details.ID)
	fmt.Println("Name:", fileDetails.Details.Name)
	fmt.Println("Description:", fileDetails.Details.Description)
}

func Get(args []string) {
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	id := getCmd.Int64("id", 0, "ID of the file")
	out := getCmd.String("o", "", "Output file path")

	err := getCmd.Parse(args)
	if err != nil || *id == 0 {
		GetUsage()
	}

	file, err := DB.GetFile(*id)
	if err != nil {
		fmt.Println("File not found")
		return
	}
	fmt.Println("File Details:")
	fmt.Println("ID:", file.Details.ID)
	fmt.Println("Name:", file.Details.Name)
	fmt.Println("Description:", file.Details.Description)

	decryptFile(file.Details.ID, *out)
}

func Update(args []string) {
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	id := getCmd.Int64("id", 0, "ID of the file")
	name := getCmd.String("n", "", "Name of the env file")
	description := getCmd.String("d", "", "Description of the env file")
	newfile := getCmd.String("f", "", "Path of env file")

	err := getCmd.Parse(args)
	if err != nil {
		log.Fatal(err)
	}
	if id == nil || *id == 0 {
		log.Fatal("Error: -id is required")
	}
	file, err := DB.GetFile(*id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Current Details:", file)
	fmt.Println("New details: ", *name, *description, *newfile)
	if *name != "" {
		file.Details.Name = *name
	}
	if *description != "" {
		file.Details.Description = *description
	}
	fmt.Println("New details: ", file)
	if *newfile == "" {
		err = DB.UpdateFile(&FileMeta{
			ID:          file.Details.ID,
			Name:        file.Details.Name,
			Description: file.Details.Description,
		})
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	encFilePath := path.Join(ENC_DIR, file.Path)
	key, salt := encryptFile(*newfile, encFilePath)
	file.Key = hex.EncodeToString(key)
	file.Salt = hex.EncodeToString(salt)
	err = DB.UpdateFileWithEncFile(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("File updated successfully")
}

func Delete(args []string) {
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	id := getCmd.Int64("id", 0, "ID of the file")

	err := getCmd.Parse(args)
	if err != nil || id == nil {
		log.Fatal(err)
	}
	err = DB.DeleteFile(*id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("File deleted successfully")
}

func List(args []string) {
	files, err := DB.GetFiles()
	if err != nil {
		log.Fatal(err)
	}
	fileTable := util.NewTable("Files")
	fileTable.InitColumns([]string{"ID", "Name", "Description", "Created at", "Updated at"})
	for _, file := range files {
		fileTable.AddRow([]string{fmt.Sprint(file.ID), file.Name, file.Description, util.GetFormattedTime(file.CreatedAt), util.GetFormattedTime(file.UpdatedAt)})
	}
	fileTable.Print()
}

func CLI() {
	cmd, args := GetCmd()
	Control(cmd, args)
}
