package challenge

import (
	"errors"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
  "gitlab.com/ncent/monetization/services/ncent/monetization/subscription"
  "gitlab.com/ncent/monetization/services/ncent/monetization/quota"
  "gitlab.com/ncent/monetization/services/aws/eventbridge/client"
)

var _ = Describe("Create Validate Challenge Suite", func() {
	var (
		ValidateChallengeService *ValidateChallengeService
		mockEventBridgeCli             *client.MockEventBridgeService
		mockSubscriptionRepo           *subscription.MockSubscriptionRepository
		mockQuotaUsageRepo             *quota.MockQuotaUsageRepository
	)

	JustBeforeEach(func() {
		ValidateChallengeService = NewValidateChallengeService()
		mockEventBridgeCli = client.NewMockEventBridgeService()
		mockSubscriptionRepo = subscription.NewMockSubscriptionRepository()
		mockQuotaUsageRepo = quota.NewMockQuotaUsageRepository()
		ValidateChallengeService.eventBridgeSvc = mockEventBridgeCli
		ValidateChallengeService.subsRepo = mockSubscriptionRepo
		ValidateChallengeService.quotaUsageRepo = mockQuotaUsageRepo
		mockSubscriptionRepo.SubscriptionReturned = subscription.SubscriptionSample
	})

	Context("Given a new challenge created", func() {
		It("It should validate the challenge", func() {
			userUUID := "55109e64-1fa1-4749-ab4f-90cb781532a7"
			err := ValidateChallengeService.Validate(userUUID, "create")
			Expect(err).To(BeNil())
		})

		It("It should not validate the challenge quota exceeded", func() {
			mockQuotaUsageRepo.QuotaUsageReturned = 11
			userUUID := "55109e64-1fa1-4749-ab4f-90cb781532a7"
			err := ValidateChallengeService.Validate(userUUID, "create")
			Expect(err).To(Equal(ErrUsageLimitExceeded))
		})

		It("I should not validate the challenge error on get subscription", func() {
			mockSubscriptionRepo.ErrorReturned = errors.New("Subscription not found")
			userUUID := "55109e64-1fa1-4749-ab4f-90cb781532a7"
			err := ValidateChallengeService.Validate(userUUID, "create")
			Expect(err).To(Not(Equal(nil)))
		})

		It("I should not validate the challenge error on check quota usage", func() {
			mockSubscriptionRepo.ErrorReturned = errors.New("Error on check quota usage")
			userUUID := "55109e64-1fa1-4749-ab4f-90cb781532a7"
			err := ValidateChallengeService.Validate(userUUID, "create")
			Expect(err).To(Not(Equal(nil)))
		})
	})

	Context("Given a new challenge shared", func() {
		It("It should validate the challenge", func() {
			userUUID := "55109e64-1fa1-4749-ab4f-90cb781532a7"
			err := ValidateChallengeService.Validate(userUUID, "shared")
			Expect(err).To(BeNil())
		})

		It("It should not validate the challenge quota exceeded", func() {
			mockQuotaUsageRepo.QuotaUsageReturned = 11
			userUUID := "55109e64-1fa1-4749-ab4f-90cb781532a7"
			err := ValidateChallengeService.Validate(userUUID, "shared")
			Expect(err).To(Equal(ErrUsageLimitExceeded))
		})

		It("I should not validate the challenge error on get subscription", func() {
			mockSubscriptionRepo.ErrorReturned = errors.New("Subscription not found")
			userUUID := "55109e64-1fa1-4749-ab4f-90cb781532a2"
			err := ValidateChallengeService.Validate(userUUID, "shared")
			Expect(err).To(Not(Equal(nil)))
		})

		It("I should not validate the challenge error on check quota usage", func() {
			mockSubscriptionRepo.ErrorReturned = errors.New("Error on check quota usage")
			userUUID := "55109e64-1fa1-4749-ab4f-90cb781532a7"
			err := ValidateChallengeService.Validate(userUUID, "shared")
			Expect(err).To(Not(Equal(nil)))
		})
	})
})

func TestValidateService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Create Validate Challenge Suite")
}
