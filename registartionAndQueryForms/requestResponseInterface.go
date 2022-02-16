package queryform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type SendRequest interface {
	Send() (*http.Client, *http.Response, error)
}

type UnmarshalAndPrintResponse interface {
	UnmarshalPrint(*http.Client, *http.Response, error)
}

func (srg *ServiceRegReq) Send() (*http.Client, *http.Response, error) {

	// Converting the object/struct v into a JSON encoding and returns a byte code of the JSON.
	payload, err := json.MarshalIndent(srg, "", " ")
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Payload printed: ", string(payload))

	serviceRegistryURL := "http://localhost:4245/serviceregistry/register"

	// Set the HTTP POST method, url and request body
	req, err := http.NewRequest(http.MethodPost, serviceRegistryURL, bytes.NewBuffer(payload))
	if err != nil {
		log.Println(err)

	}
	fmt.Println("Request body printed: ", req.Body)

	defer req.Body.Close()
	//Set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}

	resp, err := client.Do(req)

	return client, resp, err

}

func (serviceReq *ServiceRequestForm) Send() (*http.Client, *http.Response, error) {

	// Converting the object/struct v into a JSON encoding and returns a byte code of the JSON.
	payload, err := json.MarshalIndent(serviceReq, "", " ")
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Payload printed: ", string(payload))

	serviceOrchestrationURL := "http://localhost:4000/Orc"

	// Set the HTTP POST method, url and request body
	req, err := http.NewRequest(http.MethodPost, serviceOrchestrationURL, bytes.NewBuffer(payload))
	if err != nil {
		log.Println(err)

	}
	fmt.Println("Request body printed: ", req.Body)

	defer req.Body.Close()
	//Set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}

	resp, err := client.Do(req)

	return client, resp, err
}

func (regreply *RegistrationReply) UnmarshalPrint(client *http.Client, resp *http.Response, err error) {

	if err != nil {
		log.Println(err)
	} else {
		log.Println("Response status: ", resp.Status)
		log.Println("Response header: ", resp.Header)

		body, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			log.Println(readErr)
		} else {
			log.Println("Response boyd: ", string(body))
			err := json.Unmarshal(body, regreply)
			if err != nil {
				log.Println("Unmarshal body error: ", err)
			} else {
				fmt.Println("Unmarshal body ok: ", *regreply)
			}
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

func (serviceQueryListReply *OrchestrationResponse) UnmarshalPrint(client *http.Client, resp *http.Response, err error) {

	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("Response status: ", resp.Status)
		log.Println("Response header: ", resp.Header)

		body, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			log.Println(readErr.Error())
		} else {
			log.Println("Response boyd: ", string(body))
			err := json.Unmarshal(body, serviceQueryListReply)
			if err != nil {
				log.Println("Unmarshal body error: ", err.Error())
			} else {
				fmt.Println("Unmarshal body ok: ", *serviceQueryListReply)
			}
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
