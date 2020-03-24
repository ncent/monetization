package main

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"gitlab.com/ncent/monetization/services/aws/eventbridge/client"
	"gitlab.com/ncent/monetization/services/ncent/monetization/subscription"
)

var _ = Describe("Create Subscription Suite", func() {
	JustBeforeEach(func() {
		subscriptionRepository = subscription.NewMockSubscriptionRepository()
	})

	Context("Given the event to create subscription arrived", func() {
		It("Then it will create the subscription", func() {
			event := client.EventBridgeSample(client.SUBSCRIPTION_CREATED_SAMPLE)
			err := handler(context.Background(), event)
			Expect(err).To(BeNil())
		})
	})

	Context("Given the event to create subscription arrived", func() {
		It("Then it will fail to create subscription", func() {
			event := client.EventBridgeSample(client.SUBSCRIPTION_INVALID_SAMPLE)
			err := handler(context.Background(), event)
			Expect(err).To(Equal(errInvalidDetailType))
		})
	})

	Context("Given the event to cancelled subscription arrived", func() {
		It("Then it will cancel the subscription", func() {
			event := client.EventBridgeSample(client.SUBSCRIPTION_CANCELLED_SAMPLE)
			err := handler(context.Background(), event)
			Expect(err).To(BeNil())
		})
	})
})

func TestHandlerSubscriptionCreate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Create Subscription Suite")
}
