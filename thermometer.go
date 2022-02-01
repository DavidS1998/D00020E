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

	// What to execute for various page requests
	go http.HandleFunc("/", home)
	go http.HandleFunc("/get/", getTemperature)
	//go http.HandleFunc("/sendServiceReg/", registerService)

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

	fmt.Fprintf(w, "<a href='/sendServiceReg/'>Send Request </a>")
}

/*
func registerService(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<a href='/sendServiceReg/'>Send Request </a>")

	registerServiceToSR()
}


func registerServiceToSR( /*srg r.ServiceRegReq ) {

	var regreply *RegistrationReply = &RegistrationReply{}

	srg := ServiceRegReq{
		ServiceDefinition: "aa",
		ProviderSystemVar: ProviderSystem{
			SystemName:         "bb",
			Address:            "cc",
			Port:               222,
			AuthenticationInfo: "dd",
		},
		ServiceUri:    "ee",
		EndOfValidity: "ff",
		Secure:        "gg",
		Metadata: []string{
			"metadata1",
			"metadata2",
			"metadata3",
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

	// Converting the object/struct v into a JSON encoding and returns a byte code of the JSON.
	payload, err := json.MarshalIndent(srg, "", " ")
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Payload printed: ", string(payload))

	serviceRegistryURL := "http://localhost:4245/serviceregistry/register"

	// Set the HTTP POST method, url and request body
	req, err := http.NewRequest(http.MethodPost, serviceRegistryURL, bytes.NewBuffer(payload))
	if err != nil {
		log.Println(err)

	}
	fmt.Println("Request body printed: ", req.Body)

	defer req.Body.Close()
	//Set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
	} else {
		log.Println("Response status: ", resp.Status)
		log.Println("Response header: ", resp.Header)

		body, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			log.Println(readErr)
		} else {
			log.Println("Response boyd: ", string(body))
			err := json.Unmarshal(body, regreply)
			if err != nil {
				log.Println("Unmarshal body error: ", err)
			} else {
				fmt.Println("Unmarshal body ok: ", *regreply)
			}
			// registrationReply := r.RegistrationReply{}
			// unmarshallErr := json.Unmarshal(body, registrationReply)
			// if unmarshallErr != nil {
			// 	log.Println(registrationReply)
			// }
		}

	}
	defer resp.Body.Close()

	// closing any idle-connections that were previously connected from previous requests butare now in a "keep-alive state"
	client.CloseIdleConnections()
}
*/
