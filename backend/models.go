package main

type PZEMPayload struct {
	Voltage     float64 `json:"voltage"`
	Current     float64 `json:"current"`
	Power       float64 `json:"power"`
	Energy      float64 `json:"energy"`
	Frequency   float64 `json:"frequency"`
	PowerFactor float64 `json:"power_factor"`
}

type XYMDPayload struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
}
