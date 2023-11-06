package main

import (
	"fmt"
	"net/http"
)

func (a *App) searchHandler(w http.ResponseWriter, r *http.Request) {
	// This function runs every time there's a /search action performed in the webpage
	// This will eventually be tied to search_result.html
	searchval := r.FormValue("searchtext")
	a.search(searchval)
	fmt.Println(searchval)
}

func (a *App) search(searchval string) {
	fmt.Println("Search query " + searchval) // For now, print out the value of the searchbar for later use
}
