package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/ncent/monetization/services/ncent/monetization/subscription"
)

const authorizationHeader = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwdWJsaWNLZXkiOiIwNDhhNzVjMGI5NDVlM2YwNTU0NjJjMWNkYjNiNjgyZWUwN2I0ZGY1ZDIxODM5ZDg4Nzc4NzI5ZDQ4NzBiMGU1NmVjYjFmNzdkMThkMjgxY2EwNDA1YWFhMzM1ZTkzMzA4OGUxMWEzNDI4OGZjMjAyZmE3NGRiMTk3ZmFlNmE5NjYyIiwiZW1haWwiOiJmb29iYXJAZXhhbXBsZS5vcmcifQ.BRqcYnKvB8TIQXmAhZFgfQ3-3h_ooc9MJ9dJDg7gSCg"
const jwtSecret = "mysecret"

var _ = Describe("Create Get Subscriptions Suite", func() {
	JustBeforeEach(func() {
		mock := subscription.NewMockSubscriptionRepository()
		mock.SubscriptionsReturned = []subscription.Subscription{*subscription.SubscriptionSample}
		subscriptionRepository = mock
	})

	Context("Given the request on webhook arrived", func() {
		It("Then it will create a message on event bridge", func() {
			result, err := handler(events.APIGatewayProxyRequest{
				Headers: map[string]string{
					"Authorization": authorizationHeader,
				},
				Body: `{}`,
			})

			Expect(err).To(BeNil())
			log.Println(result.Body)
		})
	})
})

func TestHandlerGetSubsSuit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Create Get Subscriptions Suite")
}
