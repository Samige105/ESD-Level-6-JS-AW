package main

import (
	"bufio"
	"fmt"
	"os"
)

// Adding notes
type note struct {
	text string // Note text
	// owner int // Note ownership
}

func (r note) new_note() string {
	// Write to database
	fmt.Println(r.text)
	return r.text
}

func Notehandler() {
	// Testing notes
	var text string
	fmt.Println("Input some text for a test note")
	input := bufio.NewScanner(os.Stdin)
	if input.Scan() {
		fmt.Println(text)
		test_note := note{input.Text()}
		note.new_note(test_note)
	}
}
