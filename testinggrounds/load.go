package main

import (
	"fmt"
)

func check_err(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Bookreaders:\n1: Aidan\n2: James")
	var input string
	fmt.Scanln(&input)
	if (input == "1") {
		Aidan_read()
	}
	if (input == "2") {
		James_read()
	}
}