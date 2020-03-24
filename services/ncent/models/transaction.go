package models

import (
	"gitlab.com/ncent/monetization/services/ncent/monetization/models"
	"gitlab.com/ncent/monetization/services/ncent/monetization/subscription"
)

type State int

const (
	Created State = iota
	Attempted
	Failed
	Completed
)

func (m State) String() string {
	return [...]string{"Created", "Attempted", "Failed", "Completed"}[m]
}

type Transaction struct {
	UUID                      string `json:"uuid,omitempty"`
	State                     string `json:"state"`
	Date                      string `json:"date"`
	Value                     string `json:"value"`
	subscription.Subscription `json:"subscription,omitempty"`
	models.PaymentMethod      `json:"payment_method,omitempty"`
}

type Transactions []*Transaction
