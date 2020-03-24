package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
	Auth0Service "gitlab.com/ncent/monetization/services/auth0/client"
)

var (
	auth0Service Auth0Service.IAuth0Client
	log          *logrus.Entry
)

func init() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	log = logger.WithFields(logrus.Fields{
		"micro-service": "monetization",
		"handler":       "aws.eventbridge.create_auth0_user",
	})
}

func handler(ctx context.Context, event events.CloudWatchEvent) error {
	log.Infof("Event with source [%s] and type [%s] arrived on subcription service",
		event.Source, event.DetailType,
	)

	var subscription struct {
		Email string `json:"email"`
	}

	err := json.Unmarshal(event.Detail, &subscription)

	if err != nil {
		log.Errorf("Failed to unmarshal subscription %v", err)
		return err
	}

	err = auth0Service.RequestToken()

	if err != nil {
		log.Errorf("Failed to request token %v", err)
		return err
	}

	userDetails := Auth0Service.UserDetails{
		Email: subscription.Email,
	}

	err = auth0Service.CreateNewUser(&userDetails)

	if err != nil {
		log.Errorf("Failed to create Auth0 user %v", err)
		return err
	}

	return nil
}

func main() {
	auth0Service = Auth0Service.NewAuth0Client()
	lambda.Start(handler)
}
