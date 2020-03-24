package subscription

import (
	"gitlab.com/ncent/monetization/services/ncent/monetization/plan"
)

var SubscriptionSample = &Subscription{
	UUID:     "55109e64-1fa1-4749-ab4f-90cb781532a7",
	UserUUID: "4fff6e28-001a-11ea-9581-00155d75a6e4",
	Start:    "2019-11-05 22:22:17 UTC",
	End:      "2020-11-05 22:22:17 UTC",
	Active:   true,
	Provider: "Stripe",
	Plan: plan.Plan{
		UUID: "4ec48742-0195-11ea-afab-00155d43e1f0",
		Name: "Plan 1",
		Criteria: plan.Criteria{
			{
				UUID:       "a18c8ba0-0195-11ea-938d-00155d43e1f0",
				Name:       "Criteria 1",
				MaxCount:   10,
				Model:      plan.FREEMIUM.String(),
				Action:     plan.CREATE_CHALLENGE.String(),
				TimePeriod: plan.WEEK.String(),
			},
			{
				UUID:       "a18c8ba0-0195-11ea-938d-00155d13e1f0",
				Name:       "Criteria 1",
				MaxCount:   10,
				Model:      plan.FREEMIUM.String(),
				Action:     plan.SHARE_CHALLENGE.String(),
				TimePeriod: plan.WEEK.String(),
			},
		},
	},
}

type MockSubscriptionRepository struct {
	SubscriptionsReturned []Subscription
	SubscriptionReturned  *Subscription
	ErrorReturned         error
}

func NewMockSubscriptionRepository() *MockSubscriptionRepository {
	return &MockSubscriptionRepository{}
}

func (m *MockSubscriptionRepository) GetByUserUUID(userUUID string) (*Subscription, error) {
	return m.SubscriptionReturned, m.ErrorReturned
}

func (m MockSubscriptionRepository) CancelSubscription(subsUUID string) error {
	return m.ErrorReturned
}

func (m *MockSubscriptionRepository) Store(subscription *Subscription) error {
	return m.ErrorReturned
}

func (m *MockSubscriptionRepository) AllByUserUUID(userUUID string) ([]Subscription, error) {
	return m.SubscriptionsReturned, m.ErrorReturned
}

func (m *MockSubscriptionRepository) AllByUserEmail(userEmail string) ([]Subscription, error) {
	return m.SubscriptionsReturned, m.ErrorReturned
}
