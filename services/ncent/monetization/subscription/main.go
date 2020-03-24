package subscription

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"gitlab.com/ncent/monetization/services"
	dynamodbHelper "gitlab.com/ncent/monetization/services/ncent/aws/dynamodb"
)

var (
	ErrMoreThanOneSubscriptionActive = errors.New("More than one subscription active found")
	ErrNoSubscriptionFound           = errors.New("No subscription found")
	subscriptionsTableName           string
	log                              *logrus.Entry
)

func init() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	log = logger.WithFields(logrus.Fields{
		"service": "service.ncent.monetization.subscription",
	})

	stage := "test"
	if stageVar, ok := os.LookupEnv("STAGE"); !ok {
		log.Warnf("Environment variable STAGE isn't present assuming %s environment", stage)
	} else {
		stage = stageVar
	}

	subscriptionsTableName = "monetization-subscriptions-" + stage
}

type ISubscriptionRepository interface {
	GetByUserUUID(userUUID string) (*Subscription, error)
	AllByUserUUID(userUUID string) ([]Subscription, error)
	AllByUserEmail(userEmail string) ([]Subscription, error)
	Store(plan *Subscription) error
	CancelSubscription(subsUUID string) error
}

type SubscriptionRepository struct {
	client dynamodbiface.DynamoDBAPI
}

func NewSubscriptionRepository() *SubscriptionRepository {
	cfg := aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}

	return &SubscriptionRepository{
		client: services.NewTxErrorAwareDynamoDBClient(session.Must(session.NewSession(&cfg))),
	}
}

func (s SubscriptionRepository) GetByUserUUID(userUUID string) (*Subscription, error) {
	log.Infof("Get subscription by user UUID [%s]", userUUID)

	exprScan := dynamodbHelper.NewExpressionScanner(s.client)

	response, err := exprScan.ScanWithExpression(
		subscriptionsTableName,
		expression.Name("user_uuid").Equal(expression.Value(userUUID)),
		expression.Name("active").Equal(expression.Value(true)),
	)

	if err != nil {
		log.Errorf("Failed to query database %s: %v", subscriptionsTableName, err)
		return nil, err
	}

	var subscriptions []Subscription
	err = dynamodbattribute.UnmarshalListOfMaps(response.Items, &subscriptions)

	if err != nil {
		log.Errorf("Failed to unmarshall %s: %v", subscriptionsTableName, err)
		return nil, err
	}

	if len(subscriptions) > 1 {
		log.Errorf("More than one subscription found for %s", userUUID)
		return nil, ErrMoreThanOneSubscriptionActive
	}

	if len(subscriptions) == 0 {
		log.Errorf("No subscription found for %s", userUUID)
		return nil, ErrNoSubscriptionFound
	}

	log.Infof("Found subscription UUID [%s]", subscriptions[0].UUID)

	return &subscriptions[0], nil
}

func (s SubscriptionRepository) AllByUserEmail(userEmail string) ([]Subscription, error) {
	log.Infof("Get subscriptions by user email [%s]", userEmail)

	exprScan := dynamodbHelper.NewExpressionScanner(s.client)

	response, err := exprScan.ScanWithExpression(
		subscriptionsTableName,
		expression.Name("email").Equal(expression.Value(userEmail)),
	)

	if err != nil {
		log.Errorf("Failed to query database %s: %v", subscriptionsTableName, err)
		return nil, err
	}

	var subscriptions []Subscription
	err = dynamodbattribute.UnmarshalListOfMaps(response.Items, &subscriptions)

	if err != nil {
		log.Errorf("Failed to unmarshall %s: %v", subscriptionsTableName, err)
		return nil, err
	}

	log.Infoln("Found subscriptions", subscriptions)

	return subscriptions, nil
}

func (s SubscriptionRepository) AllByUserUUID(userUUID string) ([]Subscription, error) {
	log.Infof("Get subscriptions by user UUID [%s]", userUUID)

	exprScan := dynamodbHelper.NewExpressionScanner(s.client)

	response, err := exprScan.ScanWithExpression(
		subscriptionsTableName,
		expression.Name("user_uuid").Equal(expression.Value(userUUID)),
	)

	if err != nil {
		log.Errorf("Failed to query database %s: %v", subscriptionsTableName, err)
		return nil, err
	}

	var subscriptions []Subscription
	err = dynamodbattribute.UnmarshalListOfMaps(response.Items, &subscriptions)

	if err != nil {
		log.Errorf("Failed to unmarshall %s: %v", subscriptionsTableName, err)
		return nil, err
	}

	log.Infoln("Found subscriptions", subscriptions)

	return subscriptions, nil
}

func (s SubscriptionRepository) Store(subscription *Subscription) error {
	log.Infof("Store subscription [%s]", subscription.Name)

	prevSub, err := s.GetByUserUUID(subscription.UserUUID)

	if err != nil && err != ErrNoSubscriptionFound {
		log.Errorf("Failed to get previous subscription: %v", err)
		return err
	}

	if err == ErrNoSubscriptionFound {
		log.Infof("No previous subscription found")

		subscription.UUID = uuid.NewV4().String()
		subscription.Active = true

		newSubs, err := dynamodbattribute.MarshalMap(subscription)

		if err != nil {
			log.Errorf("Failed to marshal subscription %v", err)
			return err
		}

		_, err = s.client.PutItem(&dynamodb.PutItemInput{
			Item:      newSubs,
			TableName: aws.String(subscriptionsTableName),
		})

		if err != nil {
			log.Errorf("Failed to run the update %v", err)
			return err
		}
	} else {
		log.Infof("Found prev subscription [%s]", prevSub.UUID)

		key, err := dynamodbattribute.MarshalMap(struct {
			UUID string `json:"uuid"`
		}{prevSub.UUID})

		if err != nil {
			log.Errorf("Failed to marshal key %v", err)
			return err
		}

		updateSub, err := dynamodbattribute.MarshalMap(struct {
			Active bool `json:":active"`
		}{false})

		if err != nil {
			log.Errorf("Failed to marshal subscription %v", err)
			return err
		}

		subscription.UUID = uuid.NewV4().String()
		subscription.Active = true
		newSubs, err := dynamodbattribute.MarshalMap(subscription)

		if err != nil {
			log.Errorf("Failed to marshal subscription %v", err)
			return err
		}

		if err != nil {
			log.Errorf("Failed to put item on subscription %v", err)
			return err
		}

		txItems := &dynamodb.TransactWriteItemsInput{
			TransactItems: []*dynamodb.TransactWriteItem{
				{
					Update: &dynamodb.Update{
						Key:                       key,
						TableName:                 aws.String(subscriptionsTableName),
						ExpressionAttributeValues: updateSub,
						UpdateExpression:          aws.String("set active = :active"),
					},
				},
				{
					Put: &dynamodb.Put{
						Item:      newSubs,
						TableName: aws.String(subscriptionsTableName),
					},
				},
			},
		}

		_, err = s.client.TransactWriteItems(txItems)

		if err != nil {
			txErr := err.(services.TxRequestFailure)
			log.Infoln(txErr.CancellationReasons())
			log.Errorf("Failed to run the update %v", err)
			return err
		}
	}

	log.Infof("New subscription added [%s]", subscription.Name)

	return nil
}

func (s SubscriptionRepository) CancelSubscription(subsUUID string) error {
	log.Infof("Cancel subscription UUID [%s]", subsUUID)

	key, err := dynamodbattribute.MarshalMap(struct {
		UUID string `json:"uuid"`
	}{subsUUID})

	if err != nil {
		log.Errorf("Fail to cancel plan %v", err)
		return nil
	}

	updateSub, err := dynamodbattribute.MarshalMap(struct {
		Active bool `json:":active"`
	}{false})

	if err != nil {
		log.Errorf("Failed to marshal subscription %v", err)
		return err
	}

	_, err = s.client.UpdateItem(&dynamodb.UpdateItemInput{
		Key:                       key,
		TableName:                 aws.String(subscriptionsTableName),
		ExpressionAttributeValues: updateSub,
		UpdateExpression:          aws.String("set active = :active"),
	})

	if err != nil {
		log.Errorf("Failed to cancel subscription %v", err)
		return err
	}

	log.Infof("Subscription [%s] cancelled", subsUUID)

	return nil
}
