package main

import (
	"fmt"
	"os"
	"strings"
)

func check(e error) { // Funky error checking
	if e != nil && os.IsNotExist(e) { // If the folder directory doesn't exist, create it
		fmt.Println("Folder directory not found, Creating...")
		os.Mkdir("Books", 0750) // Create the directory. Note: I don't know what the 0750 means, but it makes it work so ¯\_(ツ)_/¯
	} else if e != nil { // If something else breaks, call it out here and stop the program
		panic(e)
	}
}

var booksid []int      // Give the books an ID
var booknames []string // Name of the book

func createList() {
	// Creating the lists
	var book_num = 1                  // Starting Value
	files, err := os.ReadDir("Books") // Get all files from the books directory
	check(err)
	for _, file := range files {
		booksid = append(booksid, book_num)
		bookname := strings.TrimSuffix(file.Name(), ".txt") // Cleanup the text
		booknames = append(booknames, bookname)
		book_num += 1
	}
	for i := 0; i < len(booksid); i++ {
		fmt.Println(booksid[i], ") ", booknames[i])
	}
}

func getBook() string {
	// User input stuff
	var book string
	if len(booksid) == 0 {
		fmt.Println("No books :(")
		os.Exit(0)
	} else if len(booksid) == 1 {
		fmt.Println("One book")
		fmt.Scanln(&book)
	} else {
		fmt.Println("Multiple books: 1 - ", len(booksid))
		fmt.Scanln(&book)
	}
	return book
}

func bookexists(testid string) (bool, int) {
	for _, file := range booksid {
		if testid == fmt.Sprint(file) {
			return true, file
		}
	}
	return false, -1
}

func read() {
	// Reading the text file
	var book = getBook()
	var testedBook, bookID = bookexists(book)
	if testedBook == false {
		for {
			// Loop until true or until an edgecase is hit
			if book == "quit" {
				fmt.Println("Bye!")
				os.Exit(0)
			} else if testedBook == true {
				break
			}
			fmt.Println("Invalid")
			fmt.Scanln(&book)
		}
	}
	fmt.Println("Reading: ", booknames[bookID-1])
}

func James_read() {
	fmt.Println("Sami's Book reader (Go Edition!)")
	createList()
	read()
}
