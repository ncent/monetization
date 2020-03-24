package main

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"gitlab.com/ncent/monetization/services/aws/eventbridge/client"
	"gitlab.com/ncent/monetization/services/ncent/monetization/validate/challenge"
)

var _ = Describe("Create Validate Suite", func() {
	JustBeforeEach(func() {
		validateService = challenge.NewMockValidateChallengeService()
	})

	Context("Given the event to create challenge arrived", func() {
		It("Then it will validate the challenge created", func() {
			event := client.EventBridgeSample(client.CHALLENGE_CREATED_SAMPLE)
			err := handler(context.Background(), event)
			Expect(err).To(BeNil())
		})
	})

	Context("Given the event to share challenge arrived", func() {
		It("Then it will validate the challenge shared", func() {
			event := client.EventBridgeSample(client.CHALLENGE_SHARED_SAMPLE)
			err := handler(context.Background(), event)
			Expect(err).To(BeNil())
		})
	})

	Context("Given the event challenge arrived", func() {
		It("Then it will fail to validate the challenge", func() {
			event := client.EventBridgeSample(client.CHALLENGE_INVALID_SAMPLE)
			err := handler(context.Background(), event)
			Expect(err).To(Equal(errInvalidDetailType))
		})
	})
})

func TestHandlerValidateChallenge(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Create Validate Suite")
}
