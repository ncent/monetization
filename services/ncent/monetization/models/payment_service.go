package models

type PaymentService struct {
	UUID   string `json:"uuid,omitempty"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}
