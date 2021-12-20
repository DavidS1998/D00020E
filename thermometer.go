package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	// What to execute for various page requests
	go http.HandleFunc("/", getTemperature)

	// Listens for incoming connections
	http.ListenAndServe(":8091", nil)
}

// Home page that includes a link to a subpage
func getTemperature(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, strconv.Itoa(readTemperature()))
}

// Returns a temperature
// TODO: Should be connected to a sensor
func readTemperature() int {
	// Sends a random number between 0 and 50 (for now)
	/* 	rand.Seed(time.Now().UnixNano())
	   	var randomNum = rand.Intn(50)

	   	return randomNum */

	return 28
}

// Register IP and port data to the Service Registry
/* func registerServiceToSR() {

} */
