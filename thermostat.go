/*
Run with
go run thermometer.go & go run thermostat.go & go run valve.go

Use
CURL localhost:8090
and
CURL localhost:8090/set/##
to print current status data, or to turn the servo

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

// Stored service variables
var currentTemperature = 0.0
var currentRadius = 0.0

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
	go http.HandleFunc("/set/", setValve)

	// Listens for incoming connections
	if err := http.ListenAndServe(":8090", nil); err != nil {
		panic(err)
	}

}

// Prints out thermostat data, such as current temperature and servo position
func home(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<p>Current temperature: </p>\n"+getFromService(thermometerServiceAddress, thermometerServicePort, ""))
	fmt.Fprintf(w, "<p>Current radius: </p>\n"+getFromService(valveServiceAddress, valveServicePort, "get"))
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<a href='/set/"+strconv.Itoa(30)+"'>Turn +30° </a>")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<a href='/set/"+strconv.Itoa(-30)+"'>Turn -30° </a>")
	fmt.Fprintf(w, "<br>")

	// Handy links to the other services
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<a href='http://localhost:8091/'>Thermometer </a>")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<a href='http://localhost:8092/'>Valve</a>")
}

// TODO: comment
func initClient() {

	ci = &ClientInfo{
		ClientName:   "Thermostat",
		ClientStatus: "Alive",
	}
	client = &http.Client{}

}

// Gets how much to turn the servo with, and forwards the
// formatted data as a query
// URL to get data from: localhost:8090/set/##
func setValve(w http.ResponseWriter, req *http.Request) {
	// Reads the value after /set/###
	path := strings.Split(req.URL.Path, "/")
	last := path[len(path)-1]

	// Convert to int
	num, err := strconv.Atoi(last)
	if err != nil {
		fmt.Println(err)
	}

	// PUT request for turning servo
	sendToValve(num)

	// Automatically redirect to home
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

// Sends PUT request to turn the servo in the Valve service
func sendToValve(degrees int) {

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

	// Sends the request, and waits for the response
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("Received response: ", resp.StatusCode)

	// closing any idle-connections that were previously connected from
	// previous requests but are now in a "keep-alive state"
	client.CloseIdleConnections()
}

// Sends a GET request to a service
// Will be formatted as ADDR:PORT/SUBPAGE/
func getFromService(addr string, port string, subpage string) string {
	// Tries connecting to the thermometer service
	resp, err := http.Get(addr + port + "/" + subpage + "/")
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
		currentTemperature = float64(num)
	}
	return value
}

// Requests the networking info for requested services
/* func requestServiceFromSR() {

} */

// Register IP and port data to the Service Registry
/* func registerServiceToSR() {

} */
