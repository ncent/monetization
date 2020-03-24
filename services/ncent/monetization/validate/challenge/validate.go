package challenge

import (
	"errors"
	"github.com/sirupsen/logrus"
	"os"

	"gitlab.com/ncent/monetization/services/aws/eventbridge/client"
	"gitlab.com/ncent/monetization/services/ncent/monetization/quota"
	"gitlab.com/ncent/monetization/services/ncent/monetization/subscription"
)

const (
	amountOfTimePeriod = 1
	detailType         = "challenge-created"
)

var (
	ErrUsageLimitExceeded = errors.New("Challenge Created Usage Limit Exceeded")
	eventBusName          string
	source                string
	log                   *logrus.Entry
)

func init() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	log = logger.WithFields(logrus.Fields{
		"service": "service.ncent.monetization.validate.challenge",
	})

	stage := "test"
	if stageVar, ok := os.LookupEnv("STAGE"); !ok {
		log.Warnf("Environment variable STAGE isn't present assuming %s environment", stage)
	} else {
		stage = stageVar
	}

	eventBusName = "ncent-" + stage
	source = stage + ".monetization"
}

type IValidateChallengeService interface {
	Validate(userUUID, action string) error
}

type ValidateChallengeService struct {
	eventBridgeSvc client.IEventBridgeService
	subsRepo       subscription.ISubscriptionRepository
	quotaUsageRepo quota.IQuotaUsageRepository
}

func NewValidateChallengeService() *ValidateChallengeService {
	return &ValidateChallengeService{
		eventBridgeSvc: client.NewEventBridgeService(),
		subsRepo:       subscription.NewSubscriptionRepository(),
		quotaUsageRepo: quota.NewQuotaUsageRepository(),
	}
}

func (vs ValidateChallengeService) Validate(userUUID, action string) error {
	log.Infof("Validate chalelnge UUID [%s] actions [%s]", userUUID, action)

	// TODO: Improve those function to something more OOP / S.O.L.I.D
	subs, err := vs.subsRepo.GetByUserUUID(userUUID)

	if err != nil {
		log.Errorf("Error while trying to get user %s subscription: %v", userUUID, err)
		return err
	}

	log.Infof("Subscription to validate found %s", subs.UUID)

	for _, criterion := range subs.Plan.Criteria {
		usageCount, err := vs.quotaUsageRepo.GetQuotaUsage(userUUID, action, amountOfTimePeriod, criterion.TimePeriod)

		if err != nil {
			log.Errorf("Error on check quota usage for user %s, error: %v", userUUID, err)
			return err
		}

		if usageCount > criterion.MaxCount {
			log.Errorf("Quota exceeded for user %s", userUUID)
			return ErrUsageLimitExceeded
		}
	}

	if vs.eventBridgeSvc.PutEvent(eventBusName, source, detailType, "{}"); err != nil {
		log.Errorf("Failed to put event on user %s with err %v", userUUID, err)
		return err
	}

	log.Infof("Challenge validated UUID [%s] action [%s]", userUUID, action)

	return nil
}
