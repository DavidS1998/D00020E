package main

import (
	q "VSCodeGo/services/registrationAndQueryForms"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	// First download this package from github with the "go get github.com/pborman/uuid"-command.
	// github.com/pborman/uuid (the external package) is a dependency (somthing our program depends on).
	// It then installs the "go". go mod tidy could also be used (although it differs from go get since it
	// not only adds missing dependencies but also cleans up missing dependencies).
	"github.com/pborman/uuid"
)

func main() {
	fmt.Println("Initializing thermometer system on port 8091")

	testRand := uuid.NewRandom()
	uuid := strings.Replace(testRand.String(), "-", "", -1)
	fmt.Println(uuid)

	// What to execute for various page requests
	//go http.HandleFunc("/", getTemperature)

	go http.HandleFunc("/", home)

	go http.HandleFunc("/sendServiceReg/", registerService)

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
/* func registerProviderToSR() {


}
*/

// Register service Service Registry

func home(w http.ResponseWriter, req *http.Request) {

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

	// When calling a method you have to call it from the interface-name first
	client, resp, err := srg.Send()

	regreply.UnmarshalPrint(client, resp, err)
}
