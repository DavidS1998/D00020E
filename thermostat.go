/*
Run with
go run thermometer.go & go run thermostat.go & go run valve.go

Then visit
http://localhost:8090/

Thermostat runs on port 	8090
Thermometer runs on port 	8091
Valve runs on port	 		8092
*/

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// TODO: This data should be requested from the Service Registry in the future
var thermometerServiceAddress = "http://localhost:"
var thermometerServicePort = "8091"
var valveServiceAddress = "http://localhost:"
var valveServicePort = "8092"

// Only turn the valve if the temperature differs this much
var temperatureTolerance = 3

// Stored temperature variables
var desiredTemperature = 0
var currentTemperature = 0

type ClientInfo struct {
	ClientName   string
	ClientStatus string
}
type ValveData struct {
	Degrees int
}

var (
	ci     *ClientInfo
	client *http.Client
	v      *ValveData
)

// Trying comment 3
func main() {
	fmt.Println("Initializing thermostat system on port 8090")
	initClient()

	// What to execute for various page requests
	go http.HandleFunc("/", home)
	go http.HandleFunc("/set/", setDesiredTemperature)
	go http.HandleFunc("/turnValve/", setValve)

	// Listens for incoming connections
	if err := http.ListenAndServe(":8090", nil); err != nil {
		panic(err)
	}

}

// Prints out thermostat data, such as desired and current temperature
func home(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<p>Current temperature: </p>\n"+getTempFromThermometer())
	fmt.Fprintf(w, "<p>Desired temperature: </p>\n"+strconv.Itoa(desiredTemperature))
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<a href='/set/"+strconv.Itoa(desiredTemperature+5)+"'> +5 </a>")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<a href='/set/"+strconv.Itoa(desiredTemperature-5)+"'> -5 </a>")

	// Handy links to the other services
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<a href='http://localhost:8091/'>Thermometer </a>")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<a href='http://localhost:8092/'>Valve</a>")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<a href='http://localhost:8090/turnValve/'>TurnValve</a>")

	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "Offset: "+strconv.Itoa(currentTemperatureOffset()))
}

func initClient() {

	ci = &ClientInfo{
		ClientName:   "Thermostat",
		ClientStatus: "Alive",
	}
	client = &http.Client{}

}

// Sets the desired temperatured according to URL parameters at
// localhost:8090/set/##
func setDesiredTemperature(w http.ResponseWriter, req *http.Request) {
	// Reads the value after /set/###
	path := strings.Split(req.URL.Path, "/")
	last := path[len(path)-1]

	// Convert to int
	num, err := strconv.Atoi(last)
	if err != nil {
		// Print error
		fmt.Println(err)
	} else {
		// Set temperature
		desiredTemperature = num
	}

	// Turns the valve based on temperature offset
	if currentTemperatureOffset() != 0 {
		var degreesToTurn = calculateDegreesToTurnValve(currentTemperatureOffset())
		turnValve(degreesToTurn)
	}

	// Automatically redirects to home
	http.Redirect(w, req, "/", http.StatusSeeOther)

}

// Get the difference between the ideal and current temperature
func currentTemperatureOffset() int {
	var difference = currentTemperature - desiredTemperature
	var offset = difference

	if difference < 0 {
		offset = -offset
	}

	// If the difference is small, do nothing
	if offset < temperatureTolerance {
		return 0
	} else {
		return difference
	}
}

// Calculates how many celsius translates to how many degrees to turn
func calculateDegreesToTurnValve(celsius int) int {
	// TODO: Find the correct formula somehow
	var turn = celsius * 1

	return turn
}

// func PutRequest(url string, data io.Reader) {

// }

func setValve(w http.ResponseWriter, req *http.Request) {

	turnValve(1) // Sends a put request
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

//
func turnValve(degrees int) {

	v = &ValveData{
		Degrees: degrees,
	}

	json, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	// Set the HTTP method, url and request body
	req, err := http.NewRequest(http.MethodPut, valveServiceAddress+valveServicePort+"/turn/", bytes.NewBuffer(json))
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	//Set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Print("Received response: ", resp.StatusCode)

	defer resp.Body.Close()

	// closing any idle-connections that were previously connected from previous requests butare now in a "keep-alive state"
	client.CloseIdleConnections()

}

// Scans the provided value from the thermometer service
func getTempFromThermometer() string {
	// Tries connecting to the thermometer service
	resp, err := http.Get(thermometerServiceAddress + thermometerServicePort)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Variable to store the temperature in
	var value = ""
	// Scans and prints the input
	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan() && i < 5; i++ {
		value = scanner.Text()
	}

	// Convert to int
	num, err := strconv.Atoi(value)
	if err != nil {
		// Print error
		fmt.Println(err)
	} else {
		// Set temperature
		currentTemperature = num
	}
	return value
}

// Requests the networking info for requested services
/* func requestServiceFromSR() {

} */

// Register IP and port data to the Service Registry
/* func registerServiceToSR() {

} */
