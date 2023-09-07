package main

import (
	"bufio"
	"fmt"
	"os"
)

// Note: Only one main function can exist in a directory. if you want this to work then put it into another folder
func main() {
	// Specify the file path
	filePath := "test.txt"

	// Read the content of the file
	fileContent, err := readLines(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Display the current content
	fmt.Println("Current content:")
	for i, line := range fileContent {
		fmt.Printf("%d: %s\n", i+1, line)
	}

	// Ask the user for the line number to edit
	var lineNumber int
	fmt.Print("Enter the line number to edit: ")
	_, err = fmt.Scan(&lineNumber)
	if err != nil || lineNumber < 1 || lineNumber > len(fileContent) {
		fmt.Println("Invalid line number")
		return
	}

	// Ask the user for the new line
	var newLine string
	fmt.Print("Enter the new line: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		newLine = scanner.Text()
	}

	// Update the content with the new line
	fileContent[lineNumber-1] = newLine

	// Write the modified content back to the file
	err = writeLines(filePath, fileContent)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("File edited successfully.")

	fmt.Println("Current content:")
	for i, line := range fileContent {
		fmt.Printf("%d: %s\n", i+1, line)

	}
	fmt.Printf("File edited successfully. Goodbye")

}

func readLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func writeLines(filePath string, lines []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			writer.Flush()
			return err
		}
	}
	return writer.Flush()
}
