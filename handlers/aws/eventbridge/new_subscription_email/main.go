package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	SESService "gitlab.com/ncent/monetization/services/aws/ses/client"
	"gitlab.com/ncent/monetization/services/ncent/monetization/mail/new_subscription"
)

const noReplyEmail = "no-reply@redb.ai"

var (
	sesService SESService.ISESService
	log        *logrus.Entry
)

func init() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	log = logger.WithFields(logrus.Fields{
		"micro-service": "monetization",
		"handler":       "aws.eventbridge.new_subscription_email",
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

	emailVars := new_subscription.NewSubscriptionEmailBodyVars{
		Email: subscription.Email,
	}

	emailBody, err := new_subscription.GenerateNewSubscriptionEmailBody(emailVars)

	if err != nil {
		log.Errorf("Failed to generate email body %v", err)
		return err
	}

	emailRequest := SESService.EmailRequest{
		Recipient: subscription.Email,
		Sender:    noReplyEmail,
		Html:      emailBody,
		Body:      emailBody,
		Subject:   "New Subscription",
	}

	err = sesService.SendEmail(emailRequest)

	if err != nil {
		log.Errorf("Failed to send the email %v", err)
		return err
	}

	return nil
}

func main() {
	sesService = SESService.NewSESService()
	lambda.Start(handler)
}
