package main

import (
	"fmt"
	"net/http"
	"strconv"
)

//Trying comment 3
//Trying comment 2
//Trying comment first
//new comment second
func main() {
	fmt.Println("Initializing thermometer system on port 8091")

	// What to execute for various page requests
	go http.HandleFunc("/", getTemperature)

	// Listens for incoming connections
	if err := http.ListenAndServe(":8091", nil); err != nil {
		panic(err)
	}
}

// Home page that includes a link to a subpage
func getTemperature(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, strconv.Itoa(readTemperature()))
}

// Returns a temperature
// TODO: Should be connected to a sensor
func readTemperature() int {
	return 28
}

// Register IP and port data to the Service Registry
/* func registerServiceToSR() {

} */
