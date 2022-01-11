package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type ValveData struct {
	Degrees int
}

var servoPosition = 90

func main() {
	fmt.Println("Initializing valve system on port 8092")

	go http.HandleFunc("/", home)
	go http.HandleFunc("/turn/", adjustServo)

	// Listens for incoming connections
	if err := http.ListenAndServe(":8092", nil); err != nil {
		panic(err)
	}

}

// Prints out servo position data
func home(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<p>Current position: </p>\n"+strconv.Itoa(servoPosition))
}

// TODO: Incorrect implementation of handling PUT requests. Temporary solution
func adjustServo(w http.ResponseWriter, req *http.Request) {
	// Reads the value after /turn/###

	var v ValveData

	err := json.NewDecoder(req.Body).Decode(&v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	servoPosition += v.Degrees

	// // Automatically redirects to home
	http.Redirect(w, req, "/", http.StatusSeeOther)
	return
}

// Register IP and port data to the Service Registry
/* func registerServiceToSR() {

} */
