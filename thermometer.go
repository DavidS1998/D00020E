package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Initializing thermometer system on port 8094")

	// What to execute for various page requests
	go http.HandleFunc("/", home)
	go http.HandleFunc("/get/", getTemperature)
	go http.HandleFunc("/sendServiceReg/", registerService)

	// Listens for incoming connections
	if err := http.ListenAndServe(":8094", nil); err != nil {
		panic(err)
	}
}

// Home page that includes a link to a subpage
func getTemperature(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, fmt.Sprintf("%.1f", readTemperature("28-00000dee453b")))
}

// Sends a command to a bash script that forwards the value
// argument to a Python script to turn the servo
func readTemperature(sensorID string) float64 {
	data, err := ioutil.ReadFile("/sys/bus/w1/devices/" + sensorID + "/w1_slave")
	if err != nil {
		return 9999.9
	}
	raw := string(data)
	if !strings.Contains(raw, " YES") {
		return 9999.9
	}
	i := strings.LastIndex(raw, "t=")
	if i == -1 {
		return 9999.9
	}
	temperature, err := strconv.ParseFloat(raw[i+2:len(raw)-1], 64)
	if err != nil {
		return 9999.9
	}
	return temperature / 1000.0
}

// Register service Service Registry
func home(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, fmt.Sprintf("%.1f", readTemperature("28-00000dee453b")))
	fmt.Fprintf(w, "\n<a href='/sendServiceReg/'>Send Request </a>")
}

func registerService(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<a href='/sendServiceReg/'>Send Request </a>")

	//registerServiceToSR()
}

/*
func registerServiceToSR( /*srg r.ServiceRegReq /* ) {

	var regreply *q.RegistrationReply = &q.RegistrationReply{}

	srg := &q.ServiceRegReq{
		ServiceDefinition: "aa",
		ProviderSystemVar: q.ProviderSystemReg{
			SystemName:         "bb",
			Address:            "cc",
			Port:               8091,
			AuthenticationInfo: "dd",
		},
		ServiceUri:    "ee",
		EndOfValidity: "ff",
		Secure:        "gg",
		Metadata: []string{
			"Thermometer",
			"Celsius",
			"Indoors",
			"metadata4",
		},

		Version: 33,
		Interfaces: []string{
			"Interface1",
			"Interface2",
			"Interface3",
			"Interface4",
		},
	}

	// When calling a method you have to call it from the interface-name first
	client, resp, err := srg.Send()

	regreply.UnmarshalPrint(client, resp, err)
}
*/
