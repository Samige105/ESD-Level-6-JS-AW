package main

import (
	"fmt"
	"net/http"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	// This will eventually be turned into a search function
	// TODO: Try to figure out how to actually get a variable from JS to here.
	r.ParseForm()
	searchval := r.FormValue("searchtext")
	fmt.Println(searchval)
}
