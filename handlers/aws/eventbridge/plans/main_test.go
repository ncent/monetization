package main

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"gitlab.com/ncent/monetization/services/aws/eventbridge/client"
	"gitlab.com/ncent/monetization/services/ncent/monetization/plan"
)

var _ = Describe("Create Plan Suite", func() {
	JustBeforeEach(func() {
		planRepository = plan.NewMockPlanRepository()
	})

	Context("Given the event to plans plan arrived", func() {
		It("Then it will plans the plan", func() {
			event := client.EventBridgeSample(client.PLAN_CREATED_SAMPLE)
			err := handler(context.Background(), event)
			Expect(err).To(BeNil())
		})
	})

	Context("Given the event to plans plan arrived", func() {
		It("Then it will fail to plans plan", func() {
			event := client.EventBridgeSample(client.PLAN_UPDATED_SAMPLE)
			err := handler(context.Background(), event)
			Expect(err).To(BeNil())
		})
	})
})

func TestHandlerPlanCreate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Create Plan Suite")
}
