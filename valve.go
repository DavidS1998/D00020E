package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var servoPosition = 90

func main() {
	go http.HandleFunc("/", home)
	go http.HandleFunc("/turn/", adjustServo)

	// Listens for incoming connections
	http.ListenAndServe(":8092", nil)
}

// Prints out servo position data
func home(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<p>Current position: </p>\n"+strconv.Itoa(servoPosition))
}

// TODO: Incorrect implementation of handling PUT requests. Temporary solution
func adjustServo(w http.ResponseWriter, req *http.Request) {
	// Should run a Python script


	// Reads the value after /turn/###
	path := strings.Split(req.URL.Path, "/")
	last := path[len(path)-1]

	// Convert to int
	num, err := strconv.Atoi(last)
	if err != nil {
		// Print error
		fmt.Println(err)
	} else {
		// Turn this many degrees
		servoPosition += num
	}

	// Automatically redirects to home
	http.Redirect(w, req, "/", http.StatusSeeOther)
	return
}

// Register IP and port data to the Service Registry
/* func registerServiceToSR() {

} */
