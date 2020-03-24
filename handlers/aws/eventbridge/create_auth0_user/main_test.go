package main

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	Auth0Service "gitlab.com/ncent/monetization/services/auth0/client"
	"gitlab.com/ncent/monetization/services/aws/eventbridge/client"
)

var _ = Describe("New Auth0 Email", func() {
	JustBeforeEach(func() {
		auth0Service = Auth0Service.NewMockAuth0Client()
	})

	Context("Given the Auth0 Service", func() {
		It("Then it will request to create an user", func() {
			event := client.EventBridgeSample(client.SUBSCRIPTION_CREATED_SAMPLE)
			err := handler(context.Background(), event)
			Expect(err).To(BeNil())
		})
	})
})

func TestHandlerCreateAuth0User(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "New Auth0 Email")
}
