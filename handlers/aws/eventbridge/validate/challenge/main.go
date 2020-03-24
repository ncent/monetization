package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"

	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gitlab.com/ncent/monetization/services/ncent/monetization/validate/challenge"
)

var (
	errInvalidDetailType = errors.New("Invalid detail type")
	validateService      challenge.IValidateChallengeService
	log                  *logrus.Entry
)

type Event struct {
	UUID string `json:"uuid"`
}

const (
	challengeCreatedType = "challenge-created"
	challengeSharedType  = "challenge-shared"
)

func init() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	log = logger.WithFields(logrus.Fields{
		"micro-service": "monetization",
		"handler":       "aws.eventbridge.validate.challenge",
	})
}

func handler(ctx context.Context, event events.CloudWatchEvent) error {
	log.Infof("Event with source [%s] and type [%s] arrived on validate challenge service",
		event.Source, event.DetailType,
	)

	var ev Event
	err := json.Unmarshal(event.Detail, &ev)

	if err != nil {
		log.Errorf("Failed to unmarshall event %v", err)
		return err
	}

	log.Infof("Validate UUID [%s] validate challenge service", ev.UUID)

	if event.DetailType == challengeCreatedType {
		return validateService.Validate(ev.UUID, "create")
	} else if event.DetailType == challengeSharedType {
		return validateService.Validate(ev.UUID, "share")
	} else {
		log.Warnf("Invalid detail type [%s] arrived on validate service", event.DetailType)
		return errInvalidDetailType
	}
}

func main() {
	validateService = challenge.NewValidateChallengeService()
	lambda.Start(handler)
}
