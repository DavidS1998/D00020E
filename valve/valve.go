package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/cgxeiji/servo"
)

type ValveData struct {
	Degrees int
}

// Can have a position between 0-180 degrees
var servoPosition = 90
var myServo *servo.Servo

func main() {
	fmt.Println("Initializing valve system on port 8092")

	// Turns the servo to a default position when initialized
	myServo = initServo()
	turnServo(servoPosition)

	go http.HandleFunc("/", home)
	go http.HandleFunc("/Valve/turn/", readTurnCommand)
	go http.HandleFunc("/Valve/get/", getCurrentPosition)

	// Listens for incoming connections
	if err := http.ListenAndServe(":8092", nil); err != nil {
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
}

// Used with GET requests to get current position
func getCurrentPosition(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, strconv.Itoa(servoPosition))
}

// Decodes the position data and normalizes it to a possible range (0-180)
func readTurnCommand(w http.ResponseWriter, req *http.Request) {
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

// Register IP and port data to the Service Registry
/* func registerServiceToSR() {
} */
