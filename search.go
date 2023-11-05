package main

import (
	"fmt"
	"net/http"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	searchval := r.FormValue("searchtext")
	fmt.Println(searchval)
}
