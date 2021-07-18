package main

type FormErrorResponse struct {
	FieldErrors map[string]string `json:"fieldErrors,omitempty"`
	GlobalError string            `json:"globalError,omitempty"`
}
