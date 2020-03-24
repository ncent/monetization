package main

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/ncent/monetization/services/aws/eventbridge/client"
	SESService "gitlab.com/ncent/monetization/services/aws/ses/client"
)

var _ = Describe("New Subscription Email", func() {
	JustBeforeEach(func() {
		sesService = SESService.NewMockSESService()
	})

	Context("Given the event to subscription created", func() {
		It("Then it will send subscription email", func() {
			event := client.EventBridgeSample(client.SUBSCRIPTION_CREATED_SAMPLE)
			err := handler(context.Background(), event)
			Expect(err).To(BeNil())
		})
	})
})

func TestHandlerNewSubscrioptionEmail(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "New Subscription Email")
}
