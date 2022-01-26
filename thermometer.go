package main

import (
	//r "VSCodeGo/services/thermometer/regforms"
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	fmt.Println("Initializing thermometer system on port 8091")

	// What to execute for various page requests
	go http.HandleFunc("/", getTemperature)

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

// Register service Service Registry
/* func registerServiceToSR(srg r.ServiceRegReq) {

	// Converting the object/struct v into a JSON encoding and returns a byte code of the JSON.
	payload, err := json.MarshalIndent(srg, "", " ")
	if err != nil {
		log.Println(err)
	}

	serviceRegistryURL := "http://hostname:4243/serviceregistry/register"

	// Set the HTTP POST method, url and request body
	req, err := http.NewRequest(http.MethodPost, serviceRegistryURL, bytes.NewBuffer(payload))
	if err != nil {
		log.Println(err)

	}
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
			log.Print(readErr)
		} else {
			log.Println("Response boyd: ", string(body))

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
