package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

type ValveData struct {
	Degrees int
}

var servoPosition = 90

func main() {

	cmd := exec.Command("./hello.py")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Println(cmd.Run())

	fmt.Println("Initializing valve system on port 8092")

	go http.HandleFunc("/", home)
	go http.HandleFunc("/turn/", adjustServo)

	// Listens for incoming connections
	if err := http.ListenAndServe(":8092", nil); err != nil {
		panic(err)
	}
}

// Prints out servo position data
func home(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<p>Current position: </p>\n"+strconv.Itoa(servoPosition))
}

//
func adjustServo(w http.ResponseWriter, req *http.Request) {
	var v ValveData

	err := json.NewDecoder(req.Body).Decode(&v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	servoPosition += v.Degrees
	fmt.Println("VALVE: Turned servo " + strconv.Itoa(v.Degrees) + " degrees to " + strconv.Itoa(servoPosition))
	runPythonScript(v.Degrees)

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
	fmt.Println(out)
}

// Register IP and port data to the Service Registry
/* func registerServiceToSR() {

} */
