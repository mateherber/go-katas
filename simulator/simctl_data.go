package simulator

import "encoding/json"

func UnmarshalSimulators(data []byte) (Simulators, error) {
	var r Simulators
	err := json.Unmarshal(data, &r)
	return r, err
}

type Simulators struct {
	DeviceTypes []DeviceType        `json:"devicetypes"`
	Runtimes    []Runtime           `json:"runtimes"`
	Devices     map[string][]Device `json:"devices"`
}

type Device struct {
	Udid                 string `json:"udid"`
	IsAvailable          bool   `json:"isAvailable"`
	DeviceTypeIdentifier string `json:"deviceTypeIdentifier"`
}

type DeviceType struct {
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
}

type Runtime struct {
	Identifier  string `json:"identifier"`
	Version     string `json:"version"`
	IsAvailable bool   `json:"isAvailable"`
}
