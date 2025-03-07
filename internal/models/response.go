package models

type SuccesResponse struct {
	SuccesData interface{}
}

type ErorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
