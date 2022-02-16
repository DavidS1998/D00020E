package queryform

type OrchestrationResponse struct {
	Response []Response `json:"response"`
}
type provider struct {
	ID                 int    `json:"id"`
	SystemName         string `json:"systemName"`
	Address            string `json:"address"`
	Port               int    `json:"port"`
	AuthenticationInfo string `json:"authenticationInfo"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
}
type service struct {
	ID                int    `json:"id"`
	ServiceDefinition string `json:"serviceDefinition"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
}

type Metadata struct {
	AdditionalProp1 string `json:"additionalProp1"`
	AdditionalProp2 string `json:"additionalProp2"`
	AdditionalProp3 string `json:"additionalProp3"`
}

type interfaces struct {
	ID            int    `json:"id"`
	CreatedAt     string `json:"createdAt"`
	InterfaceName string `json:"interfaceName"`
	UpdatedAt     string `json:"updatedAt"`
}
type authorizationTokens struct {
	InterfaceName1 string `json:"interfaceName1"`
	InterfaceName2 string `json:"interfaceName2"`
}

type Response struct {
	Provider            provider            `json:"provider"`
	Service             service             `json:"service"`
	ServiceURI          string              `json:"serviceUri"`
	Secure              string              `json:"secure"`
	Metadata            []string            `json:"metadata"`
	Interfaces          []interfaces        `json:"interfaces"`
	Version             int                 `json:"version"`
	AuthorizationTokens authorizationTokens `json:"authorizationTokens"`
	Warnings            []string            `json:"warnings"`
}
