package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed db/local.schema.sql
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
	Details FileMeta
	Path    string
	Key     string
	Salt    string
}

type FileMeta struct {
	ID          int64
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (db *DbCon) CreateFile(file *File) error {
	stmt, err := db.Prepare("INSERT INTO files (name, file_path, description,key,salt) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(file.Details.Name, file.Path, file.Details.Description, file.Key, file.Salt)
	if err != nil {
		return err
	}
	file.Details.ID, _ = res.LastInsertId()
	return nil
}

func (db *DbCon) GetFile(id int64) (*File, error) {
	file := new(File)
	err := db.QueryRow("SELECT id, name, file_path, description, key, salt,created_at,updated_at FROM files WHERE id = ?", id).Scan(&file.Details.ID, &file.Details.Name, &file.Path, &file.Details.Description, &file.Key, &file.Salt, &file.Details.CreatedAt, &file.Details.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (db *DbCon) GetFiles() ([]FileMeta, error) {
	rows, err := db.Query("SELECT id, name, description,created_at,updated_at FROM files")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	files := []FileMeta{}
	for rows.Next() {
		file := FileMeta{}
		err := rows.Scan(&file.ID, &file.Name, &file.Description, &file.CreatedAt, &file.UpdatedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

func (db *DbCon) UpdateFile(file *FileMeta) error {
	stmt, err := db.Prepare("UPDATE files SET name = ?, description = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(file.Name, file.Description, file.ID)
	if err != nil {
		return err
	}
	log.Println("Rows affected:", res)
	return nil
}

func (db *DbCon) UpdateFileWithEncFile(file *File) error {
	fmt.Println("Updating file with encrypted file")
	fmt.Println("File:", file)

	stmt, err := db.Prepare("UPDATE files SET name = ?, description = ?, key = ?, salt = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(file.Details.Name, file.Details.Description, file.Key, file.Salt, file.Details.ID)
	if err != nil {
		return err
	}
	id, _ := res.RowsAffected()
	log.Println("Rows affected:", id)
	return nil
}

func (db *DbCon) DeleteFile(id int64) error {
	file, err := db.GetFile(id)
	if err != nil {
		return err
	}
	stmt, err := db.Prepare("DELETE FROM files WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	err = os.Remove(path.Join(ENC_DIR, file.Path))
	if err != nil {
		log.Fatalln("Error deleting file:", err)
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
