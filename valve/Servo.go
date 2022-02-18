package main

import (
	"fmt"
	"log"

	"github.com/cgxeiji/servo"
)

func main() {
	// Use servo.Close() to close the connection of all servos and pi-blaster.
	defer servo.Close()

	// Create a new servo connected to gpio 11.
	myServo := servo.New(11)

	// (optional) Initialize the servo with your preferred values.
	myServo.SetPosition(90) // Set the initial position to 90 degrees.

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

	//move ther servo to 180 degress.
	myServo.MoveTo(180) // This is a non-blocking call.

	// MoveTo() returns a Waiter interface that can be used to move and wait on
	// the same line.
	myServo.MoveTo(0).Wait() // This is a blocking call.
}
