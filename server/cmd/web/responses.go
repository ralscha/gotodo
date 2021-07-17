package main

type DeleteResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type SaveResponse struct {
	Id          int64             `json:"id,omitempty"`
	Success     bool              `json:"success"`
	FieldErrors map[string]string `json:"fieldErrors,omitempty"`
	GlobalError string            `json:"globalError,omitempty"`
}
