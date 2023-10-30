package main

/*
To Do:
create notes functionality
create test data
create database, sorta tied to notes
Create users
Figure out database stuff
- Connect to database
- Read from database
- Write to database
*/

import (
	"fmt"
)

// Creating a new user
type user struct {
	userID int // Autoincrement, temp
	username string
	password string
}

func main() {
	fmt.Println("Hello world!")
	database := false
	// Change this to make it so the database will be built if it doesn't exist. maybe call an error if database hasn't been started or if it exists?
	if (!database) { 
		fmt.Println("Building database...")
		Build()
	}
	for true {
		var input int
		fmt.Println("MENU\n1: new note\n2: Test reading notes\n3:Login\n0: Quit")
		fmt.Scanln(&input)
		if (input == 1) {
			// Add new note
			Notehandler()
		} else {
			if (input == 0) {
				break
			}
		}
	}
	
}
