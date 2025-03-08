package main

import (
	"elk/elk/util"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/google/uuid"
)

func GetCmd() (string, []string) {
	if len(os.Args) < 2 {
		fmt.Println("Error: command is required")
		HelpExit()
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
		HelpExit()
	default:
		fmt.Println("Error: Invalid command!")
		HelpExit()
	}
}

func Create(args []string) {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	file := createCmd.String("f", "", "Path of env file")
	name := createCmd.String("n", "", "Name of the env file")
	description := createCmd.String("d", "Env file", "Description of the env file")

	err := createCmd.Parse(args)
	if err != nil {
		CreateUsageExit()
	}

	if *file == "" {
		util.PrintError("Error: -f is required")
		CreateUsageExit()
	}
	pathId := uuid.New().String() + ".enc"
	encFilepath := path.Join(ENC_DIR, pathId)
	key, salt, err := encryptFile(*file, encFilepath)
	if err != nil {
		util.PrintError("Error: couldnt encrypt file\nReason: " + err.Error())
		return
	}
	fileDetails := &File{}
	fileDetails.Details.Name = *name
	fileDetails.Path = pathId
	fileDetails.Salt = hex.EncodeToString(salt)
	fileDetails.Key = hex.EncodeToString(key)
	fileDetails.Details.Description = *description
	err = DB.CreateFile(fileDetails)
	if err != nil {
		util.PrintError("Error: couldnt create file in db")
	}
	util.PrintSuccess("File created successfully")
	fmt.Println("\nFile Details:")
	fmt.Println("\tID:", fileDetails.Details.ID)
	fmt.Println("\tName:", fileDetails.Details.Name)
	fmt.Println("\tDescription:", fileDetails.Details.Description)
}

func Get(args []string) {
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	id := getCmd.Int64("id", 0, "ID of the file")
	out := getCmd.String("o", "", "Output file path")

	err := getCmd.Parse(args)
	if err != nil || *id == 0 {
		GetUsageExit()
	}

	file, err := DB.GetFile(*id)
	if err != nil {
		util.PrintError("File not found")
		return
	}
	util.PrintSuccess("File found")
	fmt.Println("\nFile Details:")
	fmt.Println("\tID:", file.Details.ID)
	fmt.Println("\tName:", file.Details.Name)
	fmt.Println("\tDescription:", file.Details.Description)

	err = decryptFile(file.Details.ID, *out)
	if err != nil {
		util.PrintError("Error: couldnt decrypt file\nReason: " + err.Error())
		return
	}
}

func Update(args []string) {
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	id := getCmd.Int64("id", 0, "ID of the file")
	name := getCmd.String("n", "", "Name of the env file")
	description := getCmd.String("d", "", "Description of the env file")
	newfile := getCmd.String("f", "", "Path of env file")

	err := getCmd.Parse(args)
	if err != nil {
		util.PrintError("Error: couldnt parse args")
		return
	}
	if id == nil || *id == 0 {
		util.PrintError("Error: -id is required")
		return
	}
	file, err := DB.GetFile(*id)
	if err != nil {
		util.PrintError("Error: file not found")
		return
	}
	if *name != "" {
		file.Details.Name = *name
	}
	if *description != "" {
		file.Details.Description = *description
	}
	if *newfile == "" {
		err = DB.UpdateFile(&FileMeta{
			ID:          file.Details.ID,
			Name:        file.Details.Name,
			Description: file.Details.Description,
		})
		if err != nil {
			util.PrintError("Error: couldnt update file")
			return
		}
		return
	}
	encFilePath := path.Join(ENC_DIR, file.Path)
	key, salt, err := encryptFile(*newfile, encFilePath)
	if err != nil {
		util.PrintError("Error: couldnt encrypt file \n Reason: " + err.Error())
		return
	}
	file.Key = hex.EncodeToString(key)
	file.Salt = hex.EncodeToString(salt)
	err = DB.UpdateFileWithEncFile(file)
	if err != nil {
		util.PrintError("Error: couldnt update file")
		return
	}
	util.PrintSuccess("File updated successfully")
}

func Delete(args []string) {
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	id := getCmd.Int64("id", 0, "ID of the file")

	err := getCmd.Parse(args)
	if err != nil {
		util.PrintError("Error: couldnt parse args")
		return
	}
	if id == nil || *id == 0 {
		util.PrintError("Error: -id is required")
		return
	}
	err = DB.DeleteFile(*id)
	if err != nil {
		util.PrintError("Error: couldnt delete file")
		return
	}
	util.PrintSuccess("File deleted successfully")
}

func List(args []string) {
	files, err := DB.GetFiles()
	if err != nil {
		util.PrintError("Error: couldnt get files")
		return
	}
	if len(files) == 0 {
		util.PrintWarning("No files found")
		return
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
