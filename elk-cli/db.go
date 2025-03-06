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

func init() {
	var err error
	DB, err = NewConnection(path.Join(CONFIG_DIR, "elk.db"))
	if err != nil {
		log.Fatal(err)
	}
}

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
	ID          int64
	Name        string
	Path        string
	Description string
	Key         string
	Salt        string
}

type FileMin struct {
	ID          int64
	Name        string
	Description string
}

func (db *DbCon) CreateFile(file *File) (*File, error) {
	stmt, err := db.Prepare("INSERT INTO files (name, file_path, description,key,salt) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(file.Name, file.Path, file.Description, file.Key, file.Salt)
	if err != nil {
		return nil, err
	}
	file.ID, _ = res.LastInsertId()
	return file, nil
}

func (db *DbCon) GetFile(id int64) (*File, error) {
	file := new(File)
	err := db.QueryRow("SELECT id, name, file_path, description, key, salt FROM files WHERE id = ?", id).Scan(&file.ID, &file.Name, &file.Path, &file.Description, &file.Key, &file.Salt)
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
