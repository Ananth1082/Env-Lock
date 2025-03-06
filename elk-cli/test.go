package main

import "log"

func CryptoTest() {
	file := encryptFile("test/test.env", "test", "")
	if file == nil {
		log.Fatalln("Error encrypting file")
		return
	}
	decryptFile(file.ID)
}

func InitTest() {
	ElkInit()
}

func DBInitTest() {
	InitDB()
}
