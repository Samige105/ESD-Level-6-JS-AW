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

// dead function
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

// this sends the notes to the web page list
func (a *App) listHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("bozo 1")
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

	//fmt.Println("bozo 2")

	var shareUser = a.username

	SQL := "SELECT * FROM users"
	rows, err := a.db.Query(SQL)
	checkInternalServerError(err, w)

	var integers []int
	var userData User
	for rows.Next() {
		err = rows.Scan(&userData.Id, &userData.Username,
			&userData.Password, &userData.Role, &userData.NotesString)
		checkInternalServerError(err, w)
		if strings.EqualFold(userData.Username, shareUser) {
			var notesstrings []string = strings.Split(userData.NotesString, ",")
			integers = stringsToIntegers(notesstrings)
		}
	}

	//fmt.Println("bozo 5")

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

	rows, err = a.db.Query(SQL)
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
	//fmt.Println("bozo 6")
	// sends the note data into the web page
	var notes Notes
	for rows.Next() {
		err = rows.Scan(&notes.Id, &notes.Title,
			&notes.DateCreated, &notes.DateEdited, &notes.SizeBytes, &notes.DisplayContents, &notes.Contents)
		checkInternalServerError(err, w)
		notes.DisplayContents = truncateText(notes.Contents, 20)
		//fmt.Println(notes.Id, integers)
		if stringInSlice[int](notes.Id, integers) || a.role == 1 {
			data.NotesArray = append(data.NotesArray, notes)
		}
	}
	t, err := template.New("list.html").Funcs(funcMap).ParseFiles("tmpl/list.html")
	// Remakes the webpage with the new list of notes
	checkInternalServerError(err, w)
	err = t.Execute(w, data)
	checkInternalServerError(err, w)
}

// crates a new note and puts the data into the table
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
	//makes it so that the that the notes are all create with the right id
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
	//makes it so that it does the date when crated a note
	var cTime = fmt.Sprintf("%d-%02d-%02d", t.Year(), int(t.Month()), t.Day())
	var bytesizes = len([]rune(r.FormValue("Contents")))

	//A tepmariry container for notes
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
	//allows cleanup of a.db.Prepare() (prevents meromory leak)
	defer stmt.Close()

	_, err = stmt.Exec(notes.Id, notes.Title, notes.DateCreated,
		notes.DateEdited, notes.SizeBytes, notes.DisplayContents, notes.Contents)
	if err != nil {
		log.Printf("Execute query error")
		checkInternalServerError(err, w)
	}

	var shareUser = a.username

	SQL = "SELECT * FROM users"
	rows, err = a.db.Query(SQL)
	checkInternalServerError(err, w)
	//
	var userData User
	for rows.Next() {
		err = rows.Scan(&userData.Id, &userData.Username,
			&userData.Password, &userData.Role, &userData.NotesString)
		checkInternalServerError(err, w)
		if strings.EqualFold(userData.Username, shareUser) {
			var notesstrings []string = strings.Split(userData.NotesString, ",")
			var integers []int = stringsToIntegers(notesstrings)
			integers = append(integers, notes.Id)
			stmt, err := a.db.Prepare(`
				UPDATE users SET user_notes=$1
				WHERE id=$2
			`)
			//error check
			if err != nil {
				log.Fatal(err)
			}
			//allows cleanup of a.db.Prepare() (prevents meromory leak)
			defer stmt.Close()

			checkInternalServerError(err, w)
			res, err := stmt.Exec(integers, userData.Id)
			checkInternalServerError(err, w)
			_, err = res.RowsAffected()
			checkInternalServerError(err, w)
		}
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// makes it so when you edit a note it updates
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
	//changes the vaules
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
	//allows cleanup of a.db.Prepare() (prevents meromory leak)
	defer stmt.Close()

	checkInternalServerError(err, w)
	res, err := stmt.Exec(notes.Title, notes.DateCreated,
		notes.DateEdited, notes.SizeBytes, notes.DisplayContents, notes.Contents, notes.Id)
	checkInternalServerError(err, w)
	_, err = res.RowsAffected()
	checkInternalServerError(err, w)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// the abilty to make it so that you can share a note with a user
func (a *App) shareHandler(w http.ResponseWriter, r *http.Request) {
	a.isAuthenticated(w, r)
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

	var noteId int
	noteId, _ = strconv.Atoi(r.FormValue("Id"))
	var shareUser = r.FormValue("shareBox")

	SQL := "SELECT * FROM users"
	rows, err := a.db.Query(SQL)
	checkInternalServerError(err, w)

	var userData User
	for rows.Next() {
		err = rows.Scan(&userData.Id, &userData.Username,
			&userData.Password, &userData.Role, &userData.NotesString)
		checkInternalServerError(err, w)
		if strings.EqualFold(userData.Username, shareUser) {
			var notesstrings []string = strings.Split(userData.NotesString, ",")
			var integers []int = stringsToIntegers(notesstrings)
			integers = append(integers, noteId)
			stmt, err := a.db.Prepare(`
				UPDATE users SET user_notes=$1
				WHERE id=$2
			`)
			if err != nil {
				log.Fatal(err)
			}
			//allows cleanup of a.db.Prepare() (prevents meromory leak)
			defer stmt.Close()

			checkInternalServerError(err, w)
			res, err := stmt.Exec(integers, userData.Id)
			checkInternalServerError(err, w)
			_, err = res.RowsAffected()
			checkInternalServerError(err, w)
		}
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// the abilty to use the serch bar to find notes
func (a *App) searchHandler(w http.ResponseWriter, r *http.Request) {
	a.isAuthenticated(w, r)
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

	var shareUser = a.username

	SQL := "SELECT * FROM users"
	rows, err := a.db.Query(SQL)
	checkInternalServerError(err, w)

	var integers []int
	var userData User
	for rows.Next() {
		err = rows.Scan(&userData.Id, &userData.Username,
			&userData.Password, &userData.Role, &userData.NotesString)
		checkInternalServerError(err, w)
		if strings.EqualFold(userData.Username, shareUser) {
			var notesstrings []string = strings.Split(userData.NotesString, ",")
			integers = stringsToIntegers(notesstrings)
		}
	}

	//fmt.Println("search")

	var search = r.FormValue("SearchBar")

	data := pwrData{}

	var funcMap = template.FuncMap{}

	//fmt.Println(search)
	if search == "" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

	if search != data.lastSearch {
		data.NotesArray = nil

		SQL := "SELECT * FROM notes"
		rows, err := a.db.Query(SQL)
		checkInternalServerError(err, w)

		var notes Notes
		for rows.Next() {
			err = rows.Scan(&notes.Id, &notes.Title,
				&notes.DateCreated, &notes.DateEdited, &notes.SizeBytes, &notes.DisplayContents, &notes.Contents)
			checkInternalServerError(err, w)

			if strings.Contains(notes.Contents, search) || strings.Contains(notes.Title, search) {
				if stringInSlice[int](notes.Id, integers) || a.role == 1 {
					data.NotesArray = append(data.NotesArray, notes)
				}
			}
		}

		data.lastSearch = search
	}

	t, err := template.New("list.html").Funcs(funcMap).ParseFiles("tmpl/list.html")
	checkInternalServerError(err, w)
	err = t.Execute(w, data)
	checkInternalServerError(err, w)
}

// delete note in the database
func (a *App) deleteHandler(w http.ResponseWriter, r *http.Request) {
	a.isAuthenticated(w, r)
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
	var notesId, _ = strconv.Atoi(r.FormValue("Id"))
	stmt, err := a.db.Prepare("DELETE FROM notes WHERE id=$1") //tels the databse where to delete from
	if err != nil {
		log.Printf("Prepare delete error")
		checkInternalServerError(err, w)
	}
	//allows cleanup of a.db.Prepare() (prevents meromory leak)
	defer stmt.Close()
	res, err := stmt.Exec(notesId)
	checkInternalServerError(err, w)
	_, err = res.RowsAffected()
	checkInternalServerError(err, w)

	SQL := "SELECT * FROM users"
	rows, err := a.db.Query(SQL)
	checkInternalServerError(err, w)

	var userData User
	for rows.Next() {
		err = rows.Scan(&userData.Id, &userData.Username,
			&userData.Password, &userData.Role, &userData.NotesString)
		checkInternalServerError(err, w)
		var notesstrings []string = strings.Split(userData.NotesString, ",")
		var integers []int = stringsToIntegers(notesstrings)
		integers = removeAll(integers, notesId)
		stmt, err := a.db.Prepare(`
			UPDATE users SET user_notes=$1
			WHERE id=$2
		`)
		if err != nil {
			log.Fatal(err)
		}
		//allows cleanup of a.db.Prepare() (prevents meromory leak)
		defer stmt.Close()
		checkInternalServerError(err, w)
		res, err := stmt.Exec(integers, userData.Id)
		checkInternalServerError(err, w)
		_, err = res.RowsAffected()
		checkInternalServerError(err, w)
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently) // Refresh the page when note is deleted to show that it was deleted successfully

}

// goes to the Router
func (a *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	// Load the list page when user logs in
	a.isAuthenticated(w, r)
	http.Redirect(w, r, "/list", http.StatusMovedPermanently)
}
