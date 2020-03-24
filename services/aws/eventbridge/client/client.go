package client

import (
	"github.com/sirupsen/logrus"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"gitlab.com/ncent/monetization/services"
)

var log *logrus.Entry

func init() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	log = logger.WithFields(logrus.Fields{
		"service": "service.aws.eventbridge.client",
	})
}

type IEventBridgeService interface {
	PutEvent(eventBusName, source, detailType, jsonDetail string) error
}

type EventBridgeService struct {
	client services.ICloudWatchEvents
}

func NewEventBridgeService() *EventBridgeService {
	cfg := aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}

	return &EventBridgeService{
		client: cloudwatchevents.New(session.Must(session.NewSession(&cfg))),
	}
}

func (e *EventBridgeService) PutEvent(eventBusName, source, detailType, jsonDetail string) error {
	log.Infof("PutEvent bus [%s] source [%s] detail type [%s] json detail [%s]",
		eventBusName, source, detailType, jsonDetail,
	)

	putEvents := &cloudwatchevents.PutEventsInput{
		Entries: []*cloudwatchevents.PutEventsRequestEntry{
			{
				Detail:       aws.String(jsonDetail),
				DetailType:   aws.String(detailType),
				EventBusName: aws.String(eventBusName),
				Source:       aws.String(source),
			},
		},
	}

	eventsOuput, err := e.client.PutEvents(putEvents)

	if err != nil {
		log.Errorf("Failed to put event %v", err)
		return err
	}

	log.Infoln("EventBridge output", eventsOuput)

	return nil
}
