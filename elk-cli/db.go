package main

import (
	"database/sql"
	"log"
	"os"
	"path"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed db/schema.sql
var sqlStmt string
var DB *DbCon

type DbCon struct {
	*sql.DB
}

func NewConnection(dbFile string) (*DbCon, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	return &DbCon{db}, err
}

type File struct {
	ID          int
	Name        string
	Path        string
	Description string
	Key         string
	Salt        string
}

type FileMin struct {
	ID          int
	Name        string
	Description string
}

func (db *DbCon) CreateFile(file *File) error {
	stmt, err := db.Prepare("INSERT INTO files (name, path, description,key,salt) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(file.Name, file.Path, file.Description, file.Key, file.Salt)
	if err != nil {
		return err
	}
	return nil
}

func (db *DbCon) GetFile(id int) (*File, error) {
	file := new(File)
	err := db.QueryRow("SELECT id, name, path, description, key, salt FROM files WHERE id = ?", id).Scan(&file.ID, &file.Name, &file.Path, &file.Description, &file.Key, &file.Salt)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (db *DbCon) GetFiles() ([]FileMin, error) {
	rows, err := db.Query("SELECT id, name, description FROM files")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	files := []FileMin{}
	for rows.Next() {
		file := FileMin{}
		err := rows.Scan(&file.ID, &file.Name, &file.Description)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

func (db *DbCon) UpdateFile(file *FileMin) error {
	stmt, err := db.Prepare("UPDATE files SET name = ?, description = ? WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(file.Name, file.Description, file.ID)
	if err != nil {
		return err
	}
	return nil
}

func InitDB() {
	dbPath := path.Join(CONFIG_DIR, "elk.db")
	_, err := os.Create(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	DB, err = NewConnection(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("Error creating table: %q: %s\n", err, sqlStmt)
	}
}
