package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB

	saveChan  chan meme
	insert    *sql.Stmt
	existStmt *sql.Stmt
)

func persist() {
	for v := range saveChan {
		_, err := insert.Exec(v.name, v.src, v.views, v.videos, v.images, v.comments, v.created, v.lastUpdated)
		if err != nil {
			log.Printf("Failed to add %s to database: %s", v.name, err)
			continue
		}
		log.Println("Successfully added", v.name)
	}
}

func memeExists(name string) bool {
	query := fmt.Sprintf("SELECT COUNT(*) FROM memeInfo WHERE name = '%s'", strings.Replace(name, "'", "''", -1))
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Failed to check if meme exists for %s: %v", name, err)
		return false
	}
	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			log.Printf("Error checking if meme %s exists: %v", name, err)
			return false
		}
	}
	err = rows.Err()
	if err != nil {
		log.Printf("Error checking if meme %s exists: %v", name, err)
		return false
	}
	return count > 0
}

func initDB(filename string) error {
	var err error
	db, err = sql.Open("sqlite3", filename)
	if err != nil {
		return fmt.Errorf("Failed to open DB: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE memeInfo(
		name varchar(255) NOT NULL PRIMARY KEY, 
		src varchar(1024), 
		views int, 
		videos int, 
		images, int, 
		comments int, 
		created time, 
		lastUpdated time)`)

	if err != nil && err.Error() != "table memeInfo already exists" {
		return fmt.Errorf("Failed to create table: %v", err)
	}

	insert, err = db.Prepare("INSERT INTO memeInfo(name, src, views, videos, images,comments, created, lastUpdated) values(?,?,?,?,?,?,?,?)")
	if err != nil {
		return fmt.Errorf("Failed to prepare insert statement: %v", err)
	}

	return nil
}
