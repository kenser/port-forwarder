package dto

type NetworkInterface struct {
	Address        string `json:"address"`
	Desc           string `json:"desc"`
	DefaultGateway bool   `json:"default_gateway"`
}
