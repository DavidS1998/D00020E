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
	// Sends a random number between 0 and 50 (for now)
	/* 	rand.Seed(time.Now().UnixNano())
	   	var randomNum = rand.Intn(50)

	   	return randomNum */
	return 28
}

// Register IP and port data to the Service Registry
/* func registerServiceToSR() {

} */
