package main

import "strconv"

// The Fill Registration Form function fills out the form structure that is used when registering a service
func FillRegistrationForm(system System, service Service) ServiceRegReq {
	var form ServiceRegReq
	form.ServiceDefinition = service.ServiceDefinition
	form.ProviderSystemVar.SystemName = system.SystemName
	form.ProviderSystemVar.Address = system.Address
	form.ProviderSystemVar.Port = system.Port
	form.ProviderSystemVar.AuthenticationInfo = system.Authenication
	form.ServiceUri = "http://" + system.Address + ":" + strconv.Itoa(system.Port) + "/" + system.SystemName + service.Path
	form.EndOfValidity = "tomorrow"
	form.Secure = "INSECURE"
	form.Metadata = service.Metadata
	form.Version = service.Version
	form.Interfaces = system.Protocol
	return form
}

// For marshalling/unmarshalling(convert object/struct to JSON(byte data of the encoded JSON)/ converting from
// JSON(byte data of the encoded JSON) to object/struct) a service registration,
// a struct is used based on the IDD (Interaction-driven design) of the Service registry.
type ServiceRegReq struct {
	ServiceDefinition string         `json:"serviceDefinition"` // `json:"serviceDefinition` --> this is a struct tag the allows us to change the name of the variable in the outprint to "serviceDefinition" from "ServiceDefinition"
	ProviderSystemVar ProviderSystem `json:"providerSystem"`
	ServiceUri        string         `json:"serviceURI"`
	EndOfValidity     string         `json:"endOfValidity"`
	Secure            string         `json:"secure"`
	Metadata          []string       `json:"metadata"`
	Version           int            `json:"version"`
	Interfaces        []string       `json:"interface"`
}

type ProviderSystem struct {
	SystemName         string `json:"systemName"`
	Address            string `json:"address"`
	Port               int    `json:"port"`
	AuthenticationInfo string `json:"authenticationInfo"`
}

// To marshal/unmarshal a reply from a service registration, a struct is used based on the IDD (Interaction-driven design) of the Service registry.
// To marshal or unmarshal a reply from a service registration, a struct is used based on the IDD of Service Registry
type RegistrationReply struct {
	ID                int `json:"id"`
	ServiceDefinition struct {
		ID                int    `json:"id"`
		ServiceDefinition string `json:"serviceDefinition"`
		CreatedAt         string `json:"createdAt"`
		UpdatedAt         string `json:"updatedAt"`
	} `json:"serviceDefinition"`
	Provider struct {
		ID                 int    `json:"id"`
		SystemName         string `json:"systemName"`
		Address            string `json:"address"`
		Port               int    `json:"port"`
		AuthenticationInfo string `json:"authenticationInfo"`
		CreatedAt          string `json:"createdAt"`
		UpdatedAt          string `json:"updatedAt"`
	} `json:"provider"`
	ServiceURI    string   `json:"serviceUri"`
	EndOfValidity string   `json:"endOfValidity"`
	Secure        string   `json:"secure"`
	Metadata      []string `json:"metadata"`
	Version       int      `json:"version"`
	Interfaces    []struct {
		ID            int    `json:"id"`
		InterfaceName string `json:"interfaceName"`
		CreatedAt     string `json:"createdAt"`
		UpdatedAt     string `json:"updatedAt"`
	} `json:"interfaces"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
