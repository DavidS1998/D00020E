package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	q "providerConsumer/registartionAndQueryForms"
	"strconv"
	"strings"
)

const (
	systemName     string = "Thermometer"
	location       string = "Indoors"
	Celsius        string = "Celsius"
	CurrentVersion int    = 2
	Id                    = "499"
)

var (
	systemIpAddress string = ""

	TempratureServiceDefinition string = "Get temperature"
	TemperatureServiceName      string = "getTemperature"
	TemperatureServicePath      string = "/get/"
	TemperatureMetadata         []string
	TemperatureSensorID         string
)

var port = flag.Int("port", 8091, "listen to port")

func main() {

	flag.Parse()
	setLocalIP(port)

	// What to execute for various page requests
	go http.HandleFunc("/", home)
	go http.HandleFunc("/Thermometer/get/", getTemperature)
	go http.HandleFunc("/Thermometer/sendServiceReg/", registerServices)

	// Listens for incoming connections
	if err := http.ListenAndServe(fmt.Sprintf(":%s", strconv.Itoa(*port)), nil); err != nil {
		panic(err)
	}
}

// Page for manually registering service
func home(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<a href='/Thermometer/get/'>GET</a><br>")
	fmt.Fprintf(w, "<a href='/Thermometer/sendServiceReg/'>Send Request </a><br>")
}

// Home page that includes a link to a subpage
func getTemperature(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, fmt.Sprintf("%.2f", readTemperature("28-00000dee453b")))
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

// Used to find this system's networking addresses
func setLocalIP(port *int) {
	addrs, _ := net.InterfaceAddrs()

	// 0 is loopback, 1 is IPv4
	var IPv4 = addrs[1].String()
	IPv4 = strings.Split(IPv4, "/")[0]

	fmt.Printf("\n Running on local address " + IPv4 + ":" + strconv.Itoa(*port))

	systemIpAddress = IPv4
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
	system.Port = *port
	system.Authenication = ""
	system.Protocol = nil

}

func provideThermometerServiceSpecs(service *q.Service) {

	service.ServiceDefinition = TempratureServiceDefinition
	service.ServiceName = TemperatureServiceName
	service.Path = TemperatureServicePath
	service.Metadata = append(service.Metadata, TemperatureSensorID, location, Celsius, Id)
	service.Version = CurrentVersion

}
