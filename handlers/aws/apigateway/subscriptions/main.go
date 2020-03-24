package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
	"gitlab.com/ncent/monetization/services"
	"gitlab.com/ncent/monetization/services/ncent/monetization/subscription"
)

var (
	subscriptionRepository subscription.ISubscriptionRepository
	log                    *logrus.Entry
)

func init() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	log = logger.WithFields(logrus.Fields{
		"micro-service": "monetization",
		"handler":       "aws.apigateway.subscriptions",
	})
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{
		Headers:    make(map[string]string),
		StatusCode: http.StatusOK,
	}

	resp.Headers["Access-Control-Allow-Origin"] = "*"
	resp.Headers["Access-Control-Allow-Credentials"] = "true"

	token, err := services.GetTokenFromHeader(request)

	if err != nil {
		log.Errorf("Error on retrieve token header")
		resp.StatusCode = http.StatusBadRequest
		return resp, err
	}

	userEmail, err := services.GetJWTEmail(token)

	if err != nil {
		log.Errorf("Failed to get publickey %v", err)
		resp.StatusCode = http.StatusInternalServerError
		return resp, err
	}

	log.Infof("Found user with key %s", userEmail)

	subs, err := subscriptionRepository.AllByUserEmail(userEmail)

	if err != nil {
		log.Errorf("Failed to get subscription by user uuid %v", err)
		resp.StatusCode = http.StatusInternalServerError
		return resp, err
	}

	result, err := json.Marshal(subs)

	if err != nil {
		log.Errorf("Failed to marshal json return %v", err)
		resp.StatusCode = http.StatusInternalServerError
		return resp, err
	}

	resp.Body = string(result)

	log.Infof("Returning subs %s", string(result))

	return resp, nil
}

func main() {
	subscriptionRepository = subscription.NewSubscriptionRepository()
	lambda.Start(handler)
}
