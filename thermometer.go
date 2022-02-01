package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Initializing thermometer system on port 8091")

	testRand := uuid.NewRandom()
	uuid := strings.Replace(testRand.String(), "-", "", -1)
	fmt.Println(uuid)

	// What to execute for various page requests
	go http.HandleFunc("/", home)
	go http.HandleFunc("/get/", getTemperature)
	go http.HandleFunc("/sendServiceReg/", registerService)

	// Listens for incoming connections
	if err := http.ListenAndServe(":8091", nil); err != nil {
		panic(err)
	}
}

// Home page that includes a link to a subpage
func getTemperature(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, fmt.Sprintf("%.1f", readTemperature()))
}

// Sends a command to a bash script that forwards the value
// argument to a Python script to turn the servo
func readTemperature() float64 {
	// Call Python script
	out, err := exec.Command("/bin/sh", "gettemp.sh").Output()
	if err != nil {
		log.Fatal(err)
	}

	// Output from thermometer sensor
	var temperature = string([]byte(out))
	// Trim new line symbol
	temperature = strings.TrimSuffix(temperature, "\n")
	// Parse float from output
	if s, err := strconv.ParseFloat(temperature, 64); err == nil {
		return s
	} else {
		return -1
	}
}

// Register service Service Registry
func home(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, fmt.Sprintf("%.1f", readTemperature()))
	fmt.Fprintf(w, "<a href='/sendServiceReg/'>Send Request </a>")
}

func registerService(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<a href='/sendServiceReg/'>Send Request </a>")

	registerServiceToSR()
}

func registerServiceToSR( /*srg r.ServiceRegReq*/ ) {

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
