package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/icza/session"
)

type pwrData struct {
	Username   string
	NotesArray []Notes
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
	http.Redirect(w, r, "/notes", 301)
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
			&notes.DateCreated, &notes.DateEdited, &notes.SizeBytes, &notes.Contents)
		checkInternalServerError(err, w)
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
		http.Redirect(w, r, "/", 301)
	}

	var notes Notes
	notes.Title = r.FormValue("Title")
	notes.DateCreated = r.FormValue("DateCreated")
	notes.DateEdited = r.FormValue("DateEdited")
	notes.SizeBytes, _ = strconv.Atoi(r.FormValue("SizeBytes"))
	notes.Contents = r.FormValue("Contents")

	// Save to database
	stmt, err := a.db.Prepare(`
		INSERT INTO notes(note_title, date_created, date_edited, size_bytes, note_contents)
		VALUES($1, $2, $3, $4, $5)
	`)

	if err != nil {
		log.Printf("Prepare query error")
		checkInternalServerError(err, w)
	}
	_, err = stmt.Exec(notes.Title, notes.DateCreated,
		notes.DateEdited, notes.SizeBytes, notes.Contents)
	if err != nil {
		log.Printf("Execute query error")
		checkInternalServerError(err, w)
	}

	http.Redirect(w, r, "/", 301)
}

func (a *App) updateHandler(w http.ResponseWriter, r *http.Request) {
	a.isAuthenticated(w, r)
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 301)
	}

	var notes Notes
	notes.Id, _ = strconv.Atoi(r.FormValue("Id"))
	notes.Title = r.FormValue("Title")
	notes.DateCreated = r.FormValue("DateCreated")
	notes.DateEdited = r.FormValue("DateEdited")
	notes.SizeBytes, _ = strconv.Atoi(r.FormValue("SizeBytes"))
	notes.Contents = r.FormValue("Contents")
	stmt, err := a.db.Prepare(`
		UPDATE notes SET note_title=$1, date_created=$2, date_edited=$3, size_bytes=$4, note_contents=$5
		WHERE id=$6
	`)

	checkInternalServerError(err, w)
	res, err := stmt.Exec(notes.Title, notes.DateCreated,
		notes.DateEdited, notes.SizeBytes, notes.Contents, notes.Id)
	checkInternalServerError(err, w)
	_, err = res.RowsAffected()
	checkInternalServerError(err, w)
	http.Redirect(w, r, "/", 301)

}

func (a *App) deleteHandler(w http.ResponseWriter, r *http.Request) {
	a.isAuthenticated(w, r)
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 301)
	}
	var notesId, _ = strconv.ParseInt(r.FormValue("Id"), 10, 64)
	stmt, err := a.db.Prepare("DELETE FROM notes WHERE id=$1")
	checkInternalServerError(err, w)
	res, err := stmt.Exec(notesId)
	checkInternalServerError(err, w)
	_, err = res.RowsAffected()
	checkInternalServerError(err, w)
	http.Redirect(w, r, "/", 301)

}

func (a *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	a.isAuthenticated(w, r)
	http.Redirect(w, r, "/list", 301)
}
