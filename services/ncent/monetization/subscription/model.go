package subscription

import "gitlab.com/ncent/monetization/services/ncent/monetization/plan"

type Subscription struct {
	UUID      string `json:"uuid,omitempty"`
	Email     string `json:"email"`
	UserUUID  string `json:"user_uuid"`
	Start     string `json:"start"`
	End       string `json:"end"`
	Active    bool   `json:"active"`
	Provider  string `json:"provider"`
	plan.Plan `json:"plan"`
}

type Subscriptions []*Subscription
