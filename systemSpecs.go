package main

// A device that hosts systems has specific attributes
type Device struct {
	DeviceName        string            `json:"devicename"`
	NetworkInterfaces NetworkInterfaces `json:"networkname"`
}

type NetworkInterfaces struct {
	Ipv4 []string `json:"ipv4"`
	Ipv6 []string `json:"ipv6"`
}

// An Arrowhead system has specific properties
type System struct {
	SystemName    string   `json:"systemname"`
	Address       string   `json:"address"`
	Port          int      `json:"port"`
	Authenication string   `json:"authentication"`
	Protocol      []string `json:"protocol"`
}

// An Arrowhead service has specific properties
type Service struct {
	ServiceDefinition string   `json:"servicedefinition"`
	ServiceName       string   `json:"servicename"`
	Path              string   `json:"path"`
	Metadata          []string `json:"metadata"`
	Version           int      `json:"version"`
}
