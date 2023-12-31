package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/icza/session"
)

type pwrData struct {
	Username   string
	NotesArray []Notes
	lastSearch string
	//SearchNotesArray []Notes
}

func (a *App) notesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Method %s", r.Method)
	//if r.Method != "POST" {
	//	http.ServeFile(w, r, "tmpl/login.html")
	//	return
	//}
	data := pwrData{}
	var funcMap = template.FuncMap{
		"multiplication": func(n int, f int) int {
			return n * f
		},
		"addOne": func(n int) int {
			return n + 1
		},
	}
	t, err := template.New("notes.html").Funcs(funcMap).ParseFiles("tmpl/notes.html")
	checkInternalServerError(err, w)
	err = t.Execute(w, data)
	checkInternalServerError(err, w)
	http.Redirect(w, r, "/notes", http.StatusMovedPermanently)
}

func (a *App) listHandler(w http.ResponseWriter, r *http.Request) {
	a.isAuthenticated(w, r)

	//get the current username
	sess := session.Get(r)
	user := "[guest]"

	if sess != nil {
		user = sess.CAttr("username").(string)
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}

	// determine the sorting index
	params := mux.Vars(r)
	sortcol, err := strconv.Atoi(params["srt"])

	_, ok := params["srt"]
	if ok && err != nil {
		http.Redirect(w, r, "/list", http.StatusFound)
	}

	SQL := ""

	//sort the view data before sending it back to the template view
	switch sortcol {
	case 1:
		SQL = "SELECT * FROM notes ORDER by note_title"
	case 2:
		SQL = "SELECT * FROM notes ORDER by date_created"
	case 3:
		SQL = "SELECT * FROM notes ORDER by date_edited"
	case 4:
		SQL = "SELECT * FROM notes ORDER by size_bytes"
	default:
		SQL = "SELECT * FROM notes ORDER by id"
	}

	rows, err := a.db.Query(SQL)
	checkInternalServerError(err, w)
	var funcMap = template.FuncMap{
		"multiplication": func(n int, f int) int {
			return n * f
		},
		"addOne": func(n int) int {
			return n + 1
		},
	}

	data := pwrData{}
	data.Username = user

	var notes Notes
	for rows.Next() {
		err = rows.Scan(&notes.Id, &notes.Title,
			&notes.DateCreated, &notes.DateEdited, &notes.SizeBytes, &notes.DisplayContents, &notes.Contents)
		checkInternalServerError(err, w)
		notes.DisplayContents = truncateText(notes.Contents, 20)
		data.NotesArray = append(data.NotesArray, notes)
	}
	t, err := template.New("list.html").Funcs(funcMap).ParseFiles("tmpl/list.html")
	checkInternalServerError(err, w)
	err = t.Execute(w, data)
	checkInternalServerError(err, w)
}

func (a *App) createHandler(w http.ResponseWriter, r *http.Request) {
	a.isAuthenticated(w, r)
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

	// MAKE VALUES FOR FORM
	SQL := ""
	SQL = "SELECT * FROM notes"
	rows, err := a.db.Query(SQL)
	checkInternalServerError(err, w)

	var rowEnum = 0
	for rows.Next() {
		var outID int
		// Query for a value based on a single row.
		if err := a.db.QueryRow("SELECT id from notes where id = $1", rowEnum).Scan(&outID); err != nil {
			if err == sql.ErrNoRows { //!= nil { //
				//fmt.Printf("\nrow missing at %d :: %b\n", rowEnum, err)
				break
			}
		}
		rowEnum += 1
	}
	t := time.Now()
	var cTime = fmt.Sprintf("%d-%02d-%02d", t.Year(), int(t.Month()), t.Day())
	var bytesizes = len([]rune(r.FormValue("Contents")))

	var notes Notes
	notes.Id = rowEnum
	notes.Title = r.FormValue("Title")
	notes.DateCreated = cTime
	notes.DateEdited = cTime
	notes.SizeBytes = bytesizes
	notes.DisplayContents = truncateText(r.FormValue("Contents"), 20)
	notes.Contents = r.FormValue("Contents")

	// Save to database
	stmt, err := a.db.Prepare(`
		INSERT INTO notes(id, note_title, date_created, date_edited, size_bytes, note_display, note_contents)
		VALUES($1, $2, $3, $4, $5, $6, $7)
	`)
	if err != nil {
		log.Printf("Prepare query error")
		checkInternalServerError(err, w)
	}
	defer stmt.Close()

	_, err = stmt.Exec(notes.Id, notes.Title, notes.DateCreated,
		notes.DateEdited, notes.SizeBytes, notes.DisplayContents, notes.Contents)
	if err != nil {
		log.Printf("Execute query error")
		checkInternalServerError(err, w)
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (a *App) updateHandler(w http.ResponseWriter, r *http.Request) {
	// Updates the page with new information when called
	a.isAuthenticated(w, r)
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

	// MAKE VALUES FOR FORM
	t := time.Now()
	var cTime = fmt.Sprintf("%d-%02d-%02d", t.Year(), int(t.Month()), t.Day())
	var bytesizes = len([]rune(r.FormValue("Contents")))

	var notes Notes // Prepares a structure for holding new note information
	notes.Id, _ = strconv.Atoi(r.FormValue("Id"))
	notes.Title = r.FormValue("Title")
	notes.DateCreated = r.FormValue("DateCreated")
	notes.DateEdited = cTime
	notes.SizeBytes = bytesizes
	notes.DisplayContents = truncateText(r.FormValue("Contents"), 20)
	notes.Contents = r.FormValue("Contents")
	stmt, err := a.db.Prepare(`
		UPDATE notes SET note_title=$1, date_created=$2, date_edited=$3, size_bytes=$4, note_display=$5, note_contents=$6
		WHERE id=$7
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	checkInternalServerError(err, w)
	res, err := stmt.Exec(notes.Title, notes.DateCreated,
		notes.DateEdited, notes.SizeBytes, notes.DisplayContents, notes.Contents, notes.Id)
	checkInternalServerError(err, w)
	_, err = res.RowsAffected()
	checkInternalServerError(err, w)

	//####################//

	SQL := "SELECT * FROM users"
	rows, err := a.db.Query(SQL) // Finds all users in database
	checkInternalServerError(err, w)

	var userArray []User
	var userData User
	for rows.Next() {
		err = rows.Scan(&userData.Id, &userData.Username,
			&userData.Password, &userData.Role, &userData.Notes)
		checkInternalServerError(err, w)

		userArray = append(userArray, userData)
	}
	for id, user := range userArray {
		if user.Username == a.username { // Find all of the users notes
			userArray[id].Notes = append(user.Notes, notes.Id)
			stmt, err := a.db.Prepare(`
				UPDATE users SET notes=$1
				WHERE id=$2
			`)
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()

			checkInternalServerError(err, w)
			res, err := stmt.Exec(userArray[id].Notes, id)
			checkInternalServerError(err, w)
			_, err = res.RowsAffected()
			checkInternalServerError(err, w)
			break
		}
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (a *App) shareHandler(w http.ResponseWriter, r *http.Request) {
	// Sharing notes between users
	a.isAuthenticated(w, r)
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

	var noteId int
	noteId, _ = strconv.Atoi(r.FormValue("noteIdToShare")) // Get the required values
	var shareUser = r.FormValue("shareBox")

	fmt.Println(user)

	SQL := "SELECT * FROM users"
	rows, err := a.db.Query(SQL)
	checkInternalServerError(err, w)

	var userArray []User
	var userData User
	for rows.Next() { // Check all user data to make sure that the correct user is getting the note shared
		err = rows.Scan(&userData.Id, &userData.Username,
			&userData.Password, &userData.Role, &userData.Notes)
		checkInternalServerError(err, w)

		userArray = append(userArray, userData)
	}

	for id, user := range userArray {
		if user.Username == shareUser { // If the user being shared to is correct then add the noteID to them
			userArray[id].Notes = append(user.Notes, noteId)
			stmt, err := a.db.Prepare(`
				UPDATE users SET notes=$1
				WHERE id=$2
			`)
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()

			checkInternalServerError(err, w)
			res, err := stmt.Exec(userArray[id].Notes, id)
			checkInternalServerError(err, w)
			_, err = res.RowsAffected()
			checkInternalServerError(err, w)
			break
		}
	}

}

func (a *App) searchHandler(w http.ResponseWriter, r *http.Request) {
	// This function deals with finding notes with the title of whatever the user searches for
	a.isAuthenticated(w, r)
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

	//fmt.Println("search")

	var search = r.FormValue("SearchBar") // Get the value of the search bar

	data := pwrData{} // Create a structure for the search to have all required information

	var funcMap = template.FuncMap{}

	//fmt.Println(search)
	if search == "" { // If search is blank, redirect back to the main page
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

	if search != data.lastSearch {
		data.NotesArray = nil

		SQL := "SELECT * FROM notes"
		rows, err := a.db.Query(SQL) // Run SQL command to get all notes
		checkInternalServerError(err, w)

		var notes Notes
		for rows.Next() { // Go through every note and check for the user's input
			err = rows.Scan(&notes.Id, &notes.Title,
				&notes.DateCreated, &notes.DateEdited, &notes.SizeBytes, &notes.DisplayContents, &notes.Contents)
			checkInternalServerError(err, w)

			if strings.Contains(notes.Contents, search) || strings.Contains(notes.Title, search) {
				data.NotesArray = append(data.NotesArray, notes) // Finding specific notes with the user's input
			}
		}
		data.lastSearch = search // Adds to temporary memory what the last search command was
	}

	t, err := template.New("list.html").Funcs(funcMap).ParseFiles("tmpl/list.html") // Remakes the webpage with the new list of notes
	checkInternalServerError(err, w)
	err = t.Execute(w, data)
	checkInternalServerError(err, w)
}

func (a *App) deleteHandler(w http.ResponseWriter, r *http.Request) {
	// Deals with deleting notes. This is so the correct note gets deleted instead of the first note
	a.isAuthenticated(w, r)
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
	var notesId, _ = strconv.ParseInt(r.FormValue("Id"), 10, 64) // Get the ID of the note that needs to be deleted
	stmt, err := a.db.Prepare("DELETE FROM notes WHERE id=$1")
	if err != nil { // If something happens check the database connection
		log.Printf("Prepare delete error")
		checkInternalServerError(err, w)
	}
	defer stmt.Close()
	res, err := stmt.Exec(notesId)
	checkInternalServerError(err, w)
	_, err = res.RowsAffected()
	checkInternalServerError(err, w)
	http.Redirect(w, r, "/", http.StatusMovedPermanently) // Refresh the page when note is deleted to show that it was deleted successfully

}

func (a *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	// Load the list page when user logs in
	a.isAuthenticated(w, r)
	http.Redirect(w, r, "/list", http.StatusMovedPermanently)
}
