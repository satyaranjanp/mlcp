package common

import "encoding/json"

type request struct {
	VehicleType string `json:"vehicleType"`
	RegnNo string `json:"regno"`
	SlotId uint32 `json:"slotId"`
}

func ParseRequest(b []byte) (*request, error) {
	r := &request{}
	err := json.Unmarshal(b, r)
	return r, err
}
