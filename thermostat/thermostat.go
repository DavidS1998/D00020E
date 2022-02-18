package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	q "providerConsumer/registartionAndQueryForms"
	"strconv"
	"strings"
	"text/template"
)

// Configure port to run on
var runOnPort = 8090

// This data is to be requested from the Service Registry
var thermometerServiceAddress = "http://87.96.164.242:"
var thermometerServicePort = "8091"
var valveServiceAddress = "http://87.96.164.242:"
var valveServicePort = "8092"

type ClientInfo struct {
	ClientName   string
	ClientStatus string
}
type ValveData struct {
	Degrees int
}

var (
	ci               *ClientInfo
	thermostatClient *http.Client
	v                *ValveData
)

// Trying comment 3
func main() {
	fmt.Println("Initializing thermostat system on port " + strconv.Itoa(runOnPort))
	initClient()

	// What to execute for various page requests
	go http.HandleFunc("/", home)
	go http.HandleFunc("/set/", setValve)
	go http.HandleFunc("/requestServices/", requestService)

	// Listens for incoming connections
	if err := http.ListenAndServe(":"+strconv.Itoa(runOnPort)+"", nil); err != nil {
		panic(err)
	}

}

// Prints out thermostat data, such as current temperature and servo position
func home(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<p>Current temperature: </p>"+getFromService(thermometerServiceAddress, thermometerServicePort, "Thermometer/get"))

	// Variables to help present data in a clearer way (Percent, degrees of total)
	var max = 180.0
	var currentPosition = getFromService(thermometerServiceAddress, thermometerServicePort, "Thermometer/get")

	fmt.Fprintf(w, "\n<p>Current radius: </p>")

	// Parse float from output
	if s, err := strconv.ParseFloat(currentPosition, 64); err != nil {
		fmt.Println("Invalid input")
	} else {
		// Percent-based representation of the servo's position
		var percentage = ((float64(s) / max) * 100)
		fmt.Fprintf(w, fmt.Sprintf("%.0f", percentage)+"%%, ")
	}

	// Angle-based representation of the servo's position
	fmt.Fprintf(w, "\n"+getFromService(valveServiceAddress, valveServicePort, "Valve/get")+"°/180°")
	fmt.Fprintf(w, "\n<br>")
	fmt.Fprintf(w, "<a href='/set/"+strconv.Itoa(30)+"'>Turn +30° </a>")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<a href='/set/"+strconv.Itoa(-30)+"'>Turn -30° </a>")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<a href='/requestServices/'>Request service</a>")
	fmt.Fprintf(w, "<br>")

	// Handy links to the other services
	fmt.Fprintf(w, "<br>")

}

func initClient() {

	ci = &ClientInfo{
		ClientName:   "Thermostat",
		ClientStatus: "Alive",
	}
	thermostatClient = &http.Client{}

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
		return
		//panic(err)
	}

	// Set the HTTP method, url and request body
	req, err := http.NewRequest(http.MethodPut, valveServiceAddress+valveServicePort+"/turn/", bytes.NewBuffer(json))
	if err != nil {
		return
		//panic(err)
	}
	defer req.Body.Close()

	//Set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	// Sends the request, and waits for the response
	resp, err := thermostatClient.Do(req)
	if err != nil {
		return
		//panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("Received response: ", resp.StatusCode)

	// closing any idle-connections that were previously connected from
	// previous requests but are now in a "keep-alive state"
	thermostatClient.CloseIdleConnections()
}

// Sends a GET request to a service
// Will be formatted as ADDR:PORT/SUBPAGE/
func getFromService(addr string, port string, subpage string) string {
	// Tries connecting to the thermometer service
	resp, err := http.Get(addr + port + "/" + subpage + "/")
	if err != nil {
		fmt.Printf(err.Error())
		return "SERVICE UNAVAILABLE<br>"
		//panic(err)
	}
	defer resp.Body.Close()

	// Variable to store the temperature in
	var value = ""
	// Scans and prints the input
	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan() && i < 5; i++ {
		value = scanner.Text()
	}

	return value
}

func requestService(w http.ResponseWriter, r *http.Request) {

	// When we enter this URL path "/"
	if r.Method == "GET" {
		t, _ := template.ParseFiles("form.gtpl")

		// Writes the form to the object (second parameter) and writes it to w
		t.Execute(w, nil)
	} else {

		r.ParseForm()

		fmt.Println("Service:", r.Form["service"][0])

		var s *q.ServiceRequestForm = &q.ServiceRequestForm{}
		s.RequestedService.ServiceDefinitionRequirement = r.Form["service"][0]
		requestServiceFromOrchestrator(s)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func requestServiceFromOrchestrator(serviceReq *q.ServiceRequestForm) {

	var serviceQueryListReply *q.OrchestrationResponse = &q.OrchestrationResponse{}

	client, resp, err := serviceReq.Send()
	serviceQueryListReply.UnmarshalPrint(client, resp, err)

}

/* func requestServiceFromSR() {
} */
