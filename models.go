package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

// structures for tables to get data
type Notes struct {
	Id              int    `json:"id"`
	Title           string `json:"note_title"`
	DateCreated     string `json:"date_created"`
	DateEdited      string `json:"date_edited"`
	SizeBytes       int    `json:"size_bytes"`
	DisplayContents string `json:"note_display"`
	Contents        string `json:"note_contents"`
}

// user structer
type User struct {
	Id          int    `json:"id"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Role        int    `json:"role"`
	Notes       []int  `json:"user_notes"`
	NotesString string
}

func readData(fileName string) ([][]string, error) {
	f, err := os.Open(fileName)

	if err != nil {
		return [][]string{}, err
	}

	defer f.Close()

	r := csv.NewReader(f)

	// skip first line as this is the CSV header
	if _, err := r.Read(); err != nil {
		return [][]string{}, err
	}

	records, err := r.ReadAll()

	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}

// import the JSON data into a collection
func (a *App) importData() error {
	log.Printf("Creating tables...")
	// Create table as required, along with attribute constraints
	sql := `DROP TABLE IF EXISTS "notes";
	CREATE TABLE "notes" (
		id SERIAL PRIMARY KEY NOT NULL,
		note_title VARCHAR(255) NOT NULL,
		date_created DATE,
		date_edited DATE,
		size_bytes INTEGER,
		note_display VARCHAR(20) NOT NULL,
		note_contents VARCHAR(10485760) NOT NULL
	);`
	_, err := a.db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Table notes table created.")
	// Create a Table for the notes
	sql = `DROP TABLE IF EXISTS "users";
	CREATE TABLE "users" (
		id SERIAL PRIMARY KEY NOT NULL,
		username VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		role INTEGER DEFAULT 2 NOT NULL,
		user_notes INTEGER ARRAY NOT NULL
	);
	CREATE UNIQUE INDEX users_by_id ON users (id);`
	_, err = a.db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Table users created.")

	log.Printf("Inserting data...")

	//prepare the cost insert query
	stmt, err := a.db.Prepare("INSERT INTO notes VALUES($1,$2,$3,$4,$5,$6,$7)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// open the CSV file for importing in PG database
	data, err := readData("data/notes.csv")
	if err != nil {
		log.Fatal(err)
	}

	var c Notes
	// prepare the SQL for multiple inserts
	for _, data := range data {
		c.Id, _ = strconv.Atoi(data[0])
		c.Title = data[1]
		c.DateCreated = data[2]
		c.DateEdited = data[3]
		c.SizeBytes, _ = strconv.Atoi(data[4])
		c.DisplayContents = data[5]
		c.Contents = data[6]

		_, err := stmt.Exec(c.Id-1, c.Title, c.DateCreated, c.DateEdited, c.SizeBytes, c.DisplayContents, c.Contents)
		if err != nil {
			log.Fatal(err)
		}
	}

	//prepare the users insert query
	stmt, err = a.db.Prepare("INSERT INTO users VALUES($1,$2,$3,$4,$5)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// open the CSV file for importing in PG database
	data, err = readData("data/users.csv")
	if err != nil {
		log.Fatal(err)
	}

	var u User
	// prepare the SQL for multiple inserts
	for _, data := range data {
		u.Id, _ = strconv.Atoi(data[0])
		u.Username = data[1]
		u.Password = data[2]
		u.Role, _ = strconv.Atoi(data[3])
		u.Notes = append(u.Notes, 0)
		_, err := stmt.Exec(u.Id-1, u.Username, u.Password, u.Role, u.Notes)

		if err != nil {
			log.Fatal(err)
		}
	}

	// create temp file to notify data imported
	//can use database directly but this is an example
	// https://golangbyexample.com/touch-file-golang/
	file, err := os.Create("./imported")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	return err
}
