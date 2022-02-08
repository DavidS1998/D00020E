package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
)

type ValveData struct {
	Degrees int
}

// Can have a position between 0-180 degrees
var servoPosition = 90

func main() {
	fmt.Println("Initializing valve system on port 8092")

	// Turns the servo to a default position when initialized
	runPythonScript(servoPosition)

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
	runPythonScript(servoPosition)

	// Automatically redirects to home
	http.Redirect(w, req, "/", http.StatusSeeOther)
	return
}

// Sends a command to a bash script that forwards the value
// argument to a Python script to turn the servo
func runPythonScript(value int) {
	out, err := exec.Command("/bin/sh", "runscript.sh", strconv.Itoa(value)).Output()
	if err != nil {
		log.Fatal(err)
	}

	// The Python script will return the following byte array
	fmt.Println(string([]byte(out)))
}

// Register IP and port data to the Service Registry
/* func registerServiceToSR() {

} */
