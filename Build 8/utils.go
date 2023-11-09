package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
)

func checkInternalServerError(err error, w http.ResponseWriter) {
	// Global function to check for any back-end issues
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
	// Cuts 20 characters out of the note for display purposes, makes the page look a little nicer with longer notes
	if len(s) <= max {
		return s
	}
	if len(s) > max {
		max = max - 3
		s = s[:max] + "..."
		return s
	}
	return s[:max]
}

// Unused functions. Not exactly sure what they do.
func stringInSlice[T comparable](a T, array []T) bool {
	for _, b := range array {
		if b == a {
			return true
		}
	}
	return false
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
