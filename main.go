package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello world!")
	database := false
	// Change this to make it so the database will be built if it doesn't exist. maybe call an error if database hasn't been started or if it exists?
	if (database == false) { 
		fmt.Println("Building database...")
		Build()
	}
}
