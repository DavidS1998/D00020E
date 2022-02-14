package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/cgxeiji/servo"
)

type ValveData struct {
	Degrees int
}

// Can have a position between 0-180 degrees
var servoPosition = 90

func main() {
	fmt.Println("Initializing valve system on port 8092")

	// Turns the servo to a default position when initialized
	turnServo(servoPosition)

	go http.HandleFunc("/", home)
	go http.HandleFunc("/turn/", adjustServo)
	go http.HandleFunc("/get/", getPosition)

	// Listens for incoming connections
	if err := http.ListenAndServe(":8092", nil); err != nil {
		panic(err)
	}
}

// Prints out servo position data
func home(w http.ResponseWriter, req *http.Request) {
	// Calculate percentage between current and max position (180 degrees)
	var max = 180.0
	var percentage = (float64(servoPosition) / max) * 100
	fmt.Fprintf(w, "<p>Current position: </p>\n"+fmt.Sprintf("%.2f", percentage)+"%%")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, strconv.Itoa(servoPosition)+"°"+"/180°")
}

// Used with GET requests to get current position
func getPosition(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, strconv.Itoa(servoPosition))
}

// Decodes the position data and normalizes it to a possible range (0-180)
func adjustServo(w http.ResponseWriter, req *http.Request) {
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
	return
}

func initServo() {

}

func turnServo(value int) {
	// Use servo.Close() to close the connection of all servos and pi-blaster.
	defer servo.Close()

	// If you want to move the servos, make sure that pi-blaster is running.
	// For example, start pi-blaster as:
	// $ sudo pi-blaster --gpio 14 --pcm

	// Create a new servo connected to gpio 14.
	myServo := servo.New(11)
	// (optional) Initialize the servo with your preferred values.
	// myServo.Flags = servo.Normalized | servo.Centered
	myServo.MinPulse = 0.05 // Set the minimum pwm pulse width (default: 0.05).
	myServo.MaxPulse = 0.25 // Set the maximum pwm pulse width (default: 0.25).
	//myServo.SetPosition(90) // Set the initial position to 90 degrees.
	myServo.SetSpeed(0.2) // Set the speed to 20% (default: 1.0).
	// NOTE: The maximum speed of the servo is 0.19s/60degrees.
	// (optional) Set a verbose name.
	myServo.Name = "My Servo"

	// Print the information of the servo.
	fmt.Println(myServo)

	// Connect the servo to the daemon.
	err := myServo.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// (optional) Use myServo.Close() to close the connection to a specific
	// servo. You still need to close the connection to pi-blaster with
	// `servo.Close()`.
	defer myServo.Close()

	myServo.SetSpeed(0.5) // Set the speed to half. This is concurrent-safe.
	//myServo.MoveTo(180)   // This is a non-blocking call.

	/* do some work */

	//myServo.Wait() // Call Wait() to sync with the servo.

	// MoveTo() returns a Waiter interface that can be used to move and wait on
	// the same line.
	var floatValue = float64(value)
	myServo.MoveTo(floatValue).Wait() // This is a blocking call.
}

// Register IP and port data to the Service Registry
/* func registerServiceToSR() {

} */
