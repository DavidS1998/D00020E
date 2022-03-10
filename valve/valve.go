package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	q "providerConsumer/registartionAndQueryForms"
	"strconv"
	"strings"
	"time"

	"github.com/cgxeiji/servo"
)

const (
	systemName     string = "Valve"
	location       string = "Indoors"
	location2      string = "Boden"
	EntityOfValve  string = "Radians"
	CurrentVersion int    = 2
	Id             string = "7331"
)

var systemIpAddress string = ""

var port = flag.Int("port", 8092, "listen to port")

var (
	GetValveServiceDefinition string = "Get valve"
	GetValveServiceName       string = "getValve"
	GetValveServicePath       string = "/get/"
	GetValveMetadata          []string

	TurnValveServiceDefinition string = "Turn valve"
	TurnValveServiceName       string = "turnValve"
	TurnValveServicePath       string = "/turn/"
	TurnValveMetadata          []string

	ValveSensorID string
)

type ValveData struct {
	Degrees int
}

// Can have a position between 0-180 degrees
var servoPosition = 90
var myServo *servo.Servo

func main() {
	flag.Parse()
	setLocalIP(port)

	// Turns the servo to a default position when initialized
	myServo = initServo()
	turnServo(servoPosition)

	go http.HandleFunc("/", home)
	go http.HandleFunc("/Valve/turn/", readTurnCommand)
	go http.HandleFunc("/Valve/get/", getCurrentPosition)
	go http.HandleFunc("/Valve/sendServiceReg/", registerServices)

	// Listens for incoming connections
	if err := http.ListenAndServe(fmt.Sprintf(":%s", strconv.Itoa(*port)), nil); err != nil {
		panic(err)
	}
}

// Prints out user-facing servo position data
func home(w http.ResponseWriter, req *http.Request) {
	// Calculate percentage between current and max position (180 degrees)
	var max = 180.0
	var percentage = (float64(servoPosition) / max) * 100
	fmt.Fprintf(w, "<p>Current position: </p>\n"+fmt.Sprintf("%.2f", percentage)+"%%")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, strconv.Itoa(servoPosition)+"°"+"/180°")

	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<a href='/Valve/sendServiceReg/'>Send Request </a><br>")
}

// Used with GET requests to get current position
func getCurrentPosition(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, strconv.Itoa(servoPosition))
}

// Decodes the position data and normalizes it to a possible range (0-180)
func readTurnCommand(w http.ResponseWriter, req *http.Request) {
	fmt.Println("VALVE: PUT request received")
	// Decode JSON and get Degrees
	var v ValveData
	err := json.NewDecoder(req.Body).Decode(&v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Update internal position
	servoPosition += v.Degrees

	// The servo can only be in a position between 0 and 180 degrees.
	// Furthermore, the Python script responsible for turning the
	// servo can only handle positive values up to 180 (or it will crash)
	if servoPosition > 180 {
		servoPosition = 180
	} else if servoPosition < 0 {
		servoPosition = 0
	}

	// Update physical position
	fmt.Println("VALVE: Turning servo " + strconv.Itoa(v.Degrees) + " degrees to position " + strconv.Itoa(servoPosition))
	turnServo(servoPosition)

	// Automatically redirects to home
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

// Initializes the servo on GPIO-11 and connects it to the Pi-blaster daemon service
func initServo() *servo.Servo {
	newServo := servo.New(11)
	fmt.Println(newServo)
	newServo.Name = "servo_ID_1"

	// Connect the servo to the daemon.
	err := newServo.Connect()
	if err != nil {
		log.Fatal(err)
	}
	return newServo
}

// Turns the saved servo to x position
func turnServo(value int) {
	var floatValue = float64(value)

	// Blocking call
	myServo.MoveTo(floatValue).Wait()
	time.Sleep(time.Second * 1)
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

	var system *q.System = &q.System{}

	var serviceGetValve *q.Service = &q.Service{}
	var serviceTurnValve *q.Service = &q.Service{}

	provideValveSystemSpecs(system)

	provideGetValveServiceSpecs(serviceGetValve)
	provideTurnValveServiceSpecs(serviceTurnValve)

	registerServiceToSR(q.FillRegistrationForm(system, serviceGetValve))
	registerServiceToSR(q.FillRegistrationForm(system, serviceTurnValve))

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func registerServiceToSR(srg *q.ServiceRegReq) {

	var regreply *q.RegistrationReply = &q.RegistrationReply{}

	// When calling a method you have to call it from the interface-name first
	client, resp, err := srg.Send()

	regreply.UnmarshalPrint(client, resp, err)
}

func provideValveSystemSpecs(system *q.System) {

	system.SystemName = systemName
	system.Address = systemIpAddress
	system.Port = *port
	system.Authenication = ""
	system.Protocol = nil

}

func provideGetValveServiceSpecs(service *q.Service) {

	service.ServiceDefinition = GetValveServiceDefinition
	service.ServiceName = GetValveServiceName
	service.Path = GetValveServicePath
	service.Metadata = append(service.Metadata, ValveSensorID, location, EntityOfValve, Id, "")
	service.Version = CurrentVersion

}

func provideTurnValveServiceSpecs(service *q.Service) {

	service.ServiceDefinition = TurnValveServiceDefinition
	service.ServiceName = TurnValveServiceName
	service.Path = TurnValveServicePath
	service.Metadata = append(service.Metadata, ValveSensorID, location, location2, EntityOfValve, Id)
	service.Version = CurrentVersion

}
