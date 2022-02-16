package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	q "providerConsumer/registartionAndQueryForms"
	"strconv"
	"strings"
)

const (
	systemName      string = "Thermometer"
	systemIpAddress string = "87.96.164.242"
	systemPort      int    = 8091
	indoor          string = "Temperature indoors"
	Celsius         string = "Celsius"
	CurrentVersion  int    = 2
)

var (
	TempratureServiceDefinition string = "Get temperature"
	TemperatureServiceName      string = "getTemperature"
	TemperatureServicePath      string = "/get/"
	TemperatureMetadata         []string
	TemperatureSensorID         string
)

func main() {
	fmt.Println("Initializing thermometer system on port 8091")

	// What to execute for various page requests
	go http.HandleFunc("/", home)
	go http.HandleFunc("/Thermometer/get/", getTemperature)
	go http.HandleFunc("/Thermometer/sendServiceReg/", registerServices)

	// Listens for incoming connections
	if err := http.ListenAndServe(":8091", nil); err != nil {
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
	TemperatureSensorID = sensorID

	return temperature / 1000.0
}

// Register service Service Registry
func home(w http.ResponseWriter, req *http.Request) {
	//fmt.Fprintf(w, fmt.Sprintf("%.1f", readTemperature("28-00000dee453b")))
	fmt.Fprintf(w, "\n<a href='/Thermometer/sendServiceReg/'>Send registration request of services</a>")
}

func registerServices(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<a href='/Thermometer/sendServiceReg/'>Send Request </a>")

	var system *q.System = &q.System{}
	var service *q.Service = &q.Service{}

	provideThermometerSystemSpecs(system)
	provideThermometerServiceSpecs(service)

	registerServiceToSR(q.FillRegistrationForm(system, service))
}

func registerServiceToSR(srg *q.ServiceRegReq) {

	var regreply *q.RegistrationReply = &q.RegistrationReply{}

	// When calling a method you have to call it from the interface-name first
	client, resp, err := srg.Send()

	regreply.UnmarshalPrint(client, resp, err)
}

func provideThermometerSystemSpecs(system *q.System) {

	system.SystemName = systemName
	system.Address = systemIpAddress
	system.Port = systemPort
	system.Authenication = ""
	system.Protocol = nil

}

func provideThermometerServiceSpecs(service *q.Service) {

	service.ServiceDefinition = TempratureServiceDefinition
	service.ServiceName = TemperatureServiceName
	service.Path = TemperatureServicePath
	service.Metadata = append(service.Metadata, TemperatureSensorID, indoor, Celsius)
	service.Version = CurrentVersion

}
