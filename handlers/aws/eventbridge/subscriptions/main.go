package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gitlab.com/ncent/monetization/services/ncent/monetization/subscription"
)

var (
	errInvalidDetailType   = errors.New("Invalid detail type")
	subscriptionRepository subscription.ISubscriptionRepository
	log                    *logrus.Entry
)

const (
	subscriptionCreatedDetail   = "subscription-created"
	subscriptionCancelledDetail = "subscription-cancelled"
)

func init() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	log = logger.WithFields(logrus.Fields{
		"micro-service": "monetization",
		"handler":       "aws.eventbridge.subscriptions",
	})
}

func handler(ctx context.Context, event events.CloudWatchEvent) error {
	log.Infof("Event with source [%s] and type [%s] arrived on subcription service",
		event.Source, event.DetailType,
	)

	var subs *subscription.Subscription
	err := json.Unmarshal(event.Detail, &subs)

	if err != nil {
		log.Errorf("Failed to unmarshall event %v", err)
		return err
	}

	log.Infof("Subscription UUID [%s] subscription service", subs.UUID)

	if event.DetailType == subscriptionCreatedDetail {
		return subscriptionRepository.Store(subs)
	} else if event.DetailType == subscriptionCancelledDetail {
		return subscriptionRepository.CancelSubscription(subs.UUID)
	} else {
		log.Warnf("Invalid detail type [%s] arrived on subscription service", event.DetailType)
		return errInvalidDetailType
	}
}

func main() {
	subscriptionRepository = subscription.NewSubscriptionRepository()
	lambda.Start(handler)
}
