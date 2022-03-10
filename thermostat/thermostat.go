package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	c "providerConsumer/goCache"
	q "providerConsumer/registartionAndQueryForms"
	"strconv"
	"strings"
	"text/template"
	"time"
)

// Configure port to run on
var runOnPort = 8090

var nlc *c.LocalCache

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

var port = flag.Int("port", 8090, "listen to port")

// Trying comment 3
func main() {
	flag.Parse()
	initClient()

	nlc = c.NewLocalCache(time.Duration(time.Second * 3000))
	defer nlc.StopCleanup()

	// What to execute for various page requests
	go http.HandleFunc("/", home)
	go http.HandleFunc("/set/", setValve)
	go http.HandleFunc("/requestServices/", getServiceDefinition)

	fmt.Println("Running thermostat")

	// Listens for incoming connections
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		panic(err)
	}
}

// Prints out thermostat data, such as current temperature and servo position
func home(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<b>Thermostat</b>")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<br>")

	// If provider info exists for the given service "Get temperature" then the a get request is sent to that address
	// GET TEMPERATURE
	if pInfo, ok := nlc.Read("Get temperature"); ok == nil {
		fmt.Fprintf(w, "<p>Current temperature: </p>"+getFromService("http://"+pInfo.Address, strconv.Itoa(pInfo.Port), pInfo.SystemName+"/get"))
		fmt.Fprintf(w, "\n<br>")
	} else {
		fmt.Fprintf(w, "'Get temperature' not in cache")
		fmt.Fprintf(w, "\n<br>")
	}

	fmt.Fprintf(w, "\n<br>")

	// GET VALVE
	if pInfo, ok := nlc.Read("Get valve"); ok == nil {
		var currentPosition = getFromService("http://"+pInfo.Address, strconv.Itoa(pInfo.Port), pInfo.SystemName+"/get")

		// Variables to help present data in a clearer way (Percent, degrees of total)
		fmt.Fprintf(w, "<p>Current valve radius: </p>")
		// Parse float from output
		if s, err := strconv.ParseFloat(currentPosition, 64); err != nil {
			fmt.Println("Invalid input")
		} else {
			var max = 180.0
			// Percent-based representation of the servo's position
			var percentage = ((float64(s) / max) * 100)
			fmt.Fprintf(w, fmt.Sprintf("%.0f", percentage)+"%%, ")
		}

		//Angle-based representation of the servo's position
		fmt.Fprintf(w, "\n"+currentPosition+"°/180°")
		fmt.Fprintf(w, "\n<br>")

	} else {
		fmt.Fprintf(w, "'Get valve' not in cache")
		fmt.Fprintf(w, "\n<br>")
	}

	fmt.Fprintf(w, "\n<br>")

	// TURN VALVE
	if _, ok := nlc.Read("Turn valve"); ok == nil {
		fmt.Fprintf(w, "\n<br>")
		fmt.Fprintf(w, "<a href='/set/"+strconv.Itoa(30)+"'>Turn +30° </a>")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<a href='/set/"+strconv.Itoa(-30)+"'>Turn -30° </a>")
		fmt.Fprintf(w, "<br>")

	} else {
		fmt.Fprintf(w, "'Turn valve' not in cache")
		fmt.Fprintf(w, "\n<br>")
	}

	fmt.Fprintf(w, "\n<br>")
	fmt.Fprintf(w, "\n<br>")
	fmt.Fprintf(w, "<a href='/requestServices/'>Request service</a>")
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
	fmt.Println("THERMOSTAT: Request turning " + strconv.Itoa(num))
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
		fmt.Println("Error 1: Marshal")
		fmt.Println(err.Error())
		return
		//panic(err)
	}

	var tempAddress = ""
	var tempPort = ""
	var tempSystem = ""

	// Set the HTTP method, url and request body
	if pInfo, ok := nlc.Read("Turn valve"); ok == nil {
		tempAddress = pInfo.Address
		tempPort = strconv.Itoa(pInfo.Port)
		tempSystem = pInfo.SystemName
	} else {
		fmt.Println("No cache exists")
	}

	req, err := http.NewRequest(http.MethodPut, "http://"+tempAddress+tempPort+"/"+tempSystem+"/turn/", bytes.NewBuffer(json))
	if err != nil {
		fmt.Println("Error 2: Put")
		fmt.Println(err.Error())
		return
		//panic(err)
	}
	defer req.Body.Close()

	//Set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	// Sends the request, and waits for the response
	resp, err := thermostatClient.Do(req)
	if err != nil {
		fmt.Println("Error 3: Send")
		fmt.Println(err.Error())
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

func getServiceDefinition(w http.ResponseWriter, r *http.Request) {

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

		if r.Form["metadata"][0] != "" {
			fmt.Println("Sent METADATA 1")
			s.RequestedService.MetadataRequirements = append(s.RequestedService.MetadataRequirements, r.Form["metadata"][0])
		}
		if r.Form["metadata2"][0] != "" {
			fmt.Println("Sent METADATA 2")
			s.RequestedService.MetadataRequirements = append(s.RequestedService.MetadataRequirements, r.Form["metadata2"][0])
		}
		if r.Form["metadata3"][0] != "" {
			fmt.Println("Sent METADATA 3")
			s.RequestedService.MetadataRequirements = append(s.RequestedService.MetadataRequirements, r.Form["metadata3"][0])
		}

		//metadata extract
		requestServiceFromOrchestrator(s)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func requestServiceFromOrchestrator(serviceReq *q.ServiceRequestForm) {

	var serviceQueryListReply *q.OrchestrationResponse = &q.OrchestrationResponse{}

	client, resp, err := serviceReq.Send()
	serviceQueryListReply.UnmarshalPrint(client, resp, err)

	cacheSystemOfRequestedService(serviceQueryListReply)
}

func cacheSystemOfRequestedService(or *q.OrchestrationResponse) {

	for _, orcResp := range or.Response {

		p := c.ProviderInfo{
			SystemName: orcResp.Provider.SystemName,
			Address:    orcResp.Provider.Address + ":",
			Port:       orcResp.Provider.Port,
		}

		nlc.Update(orcResp.Service.ServiceDefinition, p, time.Now().Unix())
	}
}

/* func requestServiceFromSR() {
} */
