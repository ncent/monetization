package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
	"gitlab.com/ncent/monetization/services/ncent/monetization/plan"
)

var (
	errInvalidDetailType = errors.New("Invalid detail type")
	planRepository       plan.IPlanRepository
	log                  *logrus.Entry
)

const (
	planCreatedDetail = "plan-created"
	planUpdatedDetail = "plan-updated"
)

func init() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	log = logger.WithFields(logrus.Fields{
		"micro-service": "monetization",
		"handler":       "aws.eventbridge.plans",
	})
}

func handler(ctx context.Context, event events.CloudWatchEvent) error {
	log.Infof("Event with source [%s] and type [%s] arrived on plan service",
		event.Source, event.DetailType,
	)

	var plan *plan.Plan
	err := json.Unmarshal(event.Detail, &plan)

	if err != nil {
		log.Errorf("Failed to unmarshall event %v", err)
		return err
	}

	log.Infof("Plan UUID [%s] on plan service", plan.UUID)

	if event.DetailType == planCreatedDetail {
		return planRepository.Store(plan)
	} else if event.DetailType == planUpdatedDetail {
		updateExp := "set name = name, criteria = criteria"
		return planRepository.Update(plan.UUID, updateExp, plan)
	} else {
		log.Warnf("Invalid detail type [%s] arrived on plan service", event.DetailType)
		return errInvalidDetailType
	}
}

func main() {
	planRepository = plan.NewPlanRepository()
	lambda.Start(handler)
}
