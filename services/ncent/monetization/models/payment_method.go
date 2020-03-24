package models

type PaymentMethod struct {
	UUID           string `json:"uuid,omitempty"`
	Name           string `json:"name"`
	Active         bool   `json:"active"`
	PaymentService `json:"payment_service,omitempty"`
}
