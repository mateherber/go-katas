package version

import "encoding/json"

func UnmarshalSteps(data []byte) (Steps, error) {
	var r Steps
	err := json.Unmarshal(data, &r)
	return r, err
}

type Steps struct {
	Steps map[string]Step `json:"steps"`
}

type Step struct {
	LatestVersion string `json:"latest_version_number"`
}
