package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/ncent/monetization/services/aws/eventbridge/client"
)

var _ = Describe("Create E2E Webhook Suite", func() {
	JustBeforeEach(func() {
		eventBridgeService = client.NewMockEventBridgeService()
	})

	Context("Given the request on webhook arrived", func() {
		It("Then it will create a message on event bridge", func() {
			_, err := handler(events.APIGatewayProxyRequest{
				Body: `{
					"source": "prod.someevent",
					"event_bus_name": "ncent-development",
					"detail-type": "some.detail",
					"detail": "{\"some_detail\": \"ok\"}"
				}`,
			})

			Expect(err).To(BeNil())
		})
	})
})

func TestHandlerWebhookSuit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Create E2E Webhook Suite")
}
