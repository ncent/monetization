package plan

import (
	"errors"
	"github.com/sirupsen/logrus"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	planTableName string
	log           *logrus.Entry
)

func init() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	log = logger.WithFields(logrus.Fields{
		"service": "service.ncent.monetization.plan",
	})

	stage := "test"
	if stageVar, ok := os.LookupEnv("STAGE"); !ok {
		log.Warnf("Environment variable STAGE isn't present assuming %s environment", stage)
	} else {
		stage = stageVar
	}

	planTableName = "monetization-plans-" + stage
}

type IPlanRepository interface {
	GetByUUID(uuid string) (*Plan, error)
	Update(uuid, updateExpression string, plan *Plan) error
	Store(plan *Plan) error
}

type PlanRepository struct {
	client dynamodbiface.DynamoDBAPI
}

func NewPlanRepository() *PlanRepository {
	cfg := aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}

	return &PlanRepository{
		client: dynamodb.New(session.Must(session.NewSession(&cfg))),
	}
}

func (s PlanRepository) GetByUUID(uuid string) (*Plan, error) {
	log.Infof("Get plan by UUID %s", uuid)

	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String(planTableName),
		IndexName: aws.String("uuid"),
		KeyConditions: map[string]*dynamodb.Condition{
			"uuid": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(uuid),
					},
				},
			},
		},
	}

	response, err := s.client.Query(queryInput)
	if err != nil {
		log.Errorf("Failed to query database %s: %v", planTableName, err)
		return nil, err
	}

	plans := []Plan{}
	err = dynamodbattribute.UnmarshalListOfMaps(response.Items, &plans)

	if err != nil {
		log.Errorf("Failed to query unmarshall %s: %v", planTableName, err)
		return nil, err
	}

	if len(plans) == 0 {
		return nil, errors.New("No result returned on query")
	}

	log.Infof("Plan [%s] found", plans[0].UUID)

	return &plans[0], nil
}

func (s PlanRepository) Store(plan *Plan) error {
	log.Infof("Store plan %s", plan.Name)

	subs, err := dynamodbattribute.MarshalMap(plan)

	if err != nil {
		log.Errorf("Failed to marshall plan %v", err)
		return err
	}

	_, err = s.client.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(planTableName),
		Item:      subs,
	})

	if err != nil {
		log.Errorf("Failed to put item %v", err)
		return err
	}

	log.Infof("Plan %s saved with success", plan.Name)

	return nil
}

func (s PlanRepository) Update(uuid, updateExpression string, plan *Plan) error {
	log.Infof("Update plan %s", plan.UUID)

	key, err := dynamodbattribute.MarshalMap(struct {
		UUID string `json:"uuid"`
	}{uuid})

	if err != nil {
		log.Errorf("Failed to marshal key subscription %v", err)
		return err
	}

	updateSub, err := dynamodbattribute.MarshalMap(plan)

	if err != nil {
		log.Errorf("Failed to marshal subscription %v", err)
		return err
	}

	s.client.UpdateItem(&dynamodb.UpdateItemInput{
		Key:                       key,
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: updateSub,
		TableName:                 aws.String(planTableName),
	})

	if err != nil {
		log.Errorf("Failed to update plan %v", err)
		return err
	}

	log.Infof("Plan %s updated with success", uuid)

	return nil
}
