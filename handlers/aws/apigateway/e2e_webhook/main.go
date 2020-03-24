package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
	"gitlab.com/ncent/monetization/services/aws/eventbridge/client"
)

type Request struct {
	Source       string `json:"source"`
	EventBusName string `json:"event-bus-name"`
	DetailType   string `json:"detail-type"`
	Detail       string `json:"detail"`
}

var (
	eventBridgeService client.IEventBridgeService
	log                *logrus.Entry
)

func init() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	log = logger.WithFields(logrus.Fields{
		"micro-service": "monetization",
		"handler":       "aws.apigateway.e2e_webhook",
	})
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{
		Headers:    make(map[string]string),
		StatusCode: http.StatusOK,
	}

	resp.Headers["Access-Control-Allow-Origin"] = "*"

	var req Request
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		log.Errorf("Failed to unmarshal request body %v", err)
		resp.StatusCode = http.StatusInternalServerError
		return resp, err
	}

	err = eventBridgeService.PutEvent(
		req.EventBusName,
		req.Source,
		req.DetailType,
		req.Detail,
	)

	if err != nil {
		log.Errorf("Failed to put event %v", err)
		resp.StatusCode = http.StatusInternalServerError
		return resp, err
	}

	resp.Body = request.Body
	return resp, nil
}

func main() {
	eventBridgeService = client.NewEventBridgeService()
	lambda.Start(handler)
}
