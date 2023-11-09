package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func checkInternalServerError(err error, w http.ResponseWriter) {
	// Global function to check for any back-end issues
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// used to auto detect the active local IP address - not used yet
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func truncateText(s string, max int) string {
	if len(s) <= max {
		return s
	}
	// Cuts 20 characters out of the note for display purposes, makes the page look a little nicer with longer notes
	if len(s) > max {
		max = max - 3
		s = s[:max] + "..."
		return s
	}
	return s[:max]
}

func stringInSlice[T comparable](a T, array []T) bool {
	for _, b := range array {
		if b == a {
			return true
		}
	}
	return false
}

func stringsToIntegers(lines []string) []int {
	var integers []int
	for _, line := range lines {
		//fmt.Println(i, line)
		line = strings.ReplaceAll(line, "{", "")
		line = strings.ReplaceAll(line, "}", "")

		n, err := strconv.Atoi(line)
		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(2)
			continue
		}
		integers = append(integers, n)
	}

	return integers
}
func removeAll(nums []int, val int) []int {
	lenArr := len(nums)
	var k int = 0
	for i := 0; i < lenArr; {
		if nums[i] != val {
			nums[k] = nums[i]
			k++
		}
		i++
	}
	return nums[0:k]
}
