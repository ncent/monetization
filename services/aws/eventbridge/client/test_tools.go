package client

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
)

type MockCloudWatchEvents struct{}

func NewMockCloudWatchEvents() *MockCloudWatchEvents {
	return &MockCloudWatchEvents{}
}

func (m *MockCloudWatchEvents) PutEvents(putEventsInput *cloudwatchevents.PutEventsInput) (*cloudwatchevents.PutEventsOutput, error) {
	return nil, nil
}

type MockEventBridgeService struct {
	ErrorReturned error
}

func NewMockEventBridgeService() *MockEventBridgeService {
	return &MockEventBridgeService{}
}

func (m *MockEventBridgeService) PutEvent(eventBusName, source, detailType, jsonDetail string) error {
	return m.ErrorReturned
}

const challengeCreatedJSON string = `{
	"version": "0",
	"id": "1234",
	"detail-type": "challenge-created",
	"source": "prod.user.challenge",
	"account": "4321",
	"time": "2019-12-22T18:43:48Z",
	"region": "us-east-1",
	"resources": [],
	"detail": {
		"uuid": "99a7630f-b3c4-403c-9268-bb7dc1451e14"
	}
}`

const challengeSharedJSON string = `{
	"version": "0",
	"id": "1234",
	"detail-type": "challenge-shared",
	"source": "prod.user.challenge",
	"account": "4321",
	"time": "2019-12-22T18:43:48Z",
	"region": "us-east-1",
	"resources": [],
	"detail": {
		"uuid": "99a7630f-b3c4-403c-9268-bb7dc1451e14"
	}
}`

const challengeInvalidJSON string = `{
	"version": "0",
	"id": "1234",
	"detail-type": "challenge-dontknow",
	"source": "prod.user.challenge",
	"account": "4321",
	"time": "2019-12-22T18:43:48Z",
	"region": "us-east-1",
	"resources": [],
	"detail": {}
}`

const subscriptionCreatedJSON string = `{
	"version": "0",
	"id": "1234",
	"detail-type": "subscription-created",
	"source": "prod.user.subscription",
	"account": "4321",
	"time": "2019-12-22T18:43:48Z",
	"region": "us-east-1",
	"resources": [],
	"detail": {
		"uuid": "99a7630f-b3c4-403c-9268-bb7dc1451e14",
		"start": "2019-12-22T18:43:48Z",
		"end": "2020-12-22T18:43:48Z",
		"email": "foobartest@example.org",
		"provider": "stripe",
		"plan": {
			"uuid": "88bd0734-04cf-11ea-b52b-00155ddc81f9",
			"name": "Plan 1",
			"criteria": [{
				"uuid": "950470f4-04cf-11ea-8a5a-00155ddc81f8",
				"name": "Criteria 1",
				"max_count": 10,
				"model": "FREEMIUM",
				"action": "CREATE_CHALLENGE",
				"time_period": "WEEK"
			}]
		}
	}
}`

const subscriptionCancelledJSON string = `{
	"version": "0",
	"id": "1234",
	"detail-type": "subscription-created",
	"source": "prod.user.subscription",
	"account": "4321",
	"time": "2019-12-22T18:43:48Z",
	"region": "us-east-1",
	"resources": [],
	"detail": {
		"uuid": "99a7630f-b3c4-403c-9268-bb7dc1451e14"
	}
}`

const subscriptionInvalidJSON string = `{
	"version": "0",
	"id": "1234",
	"detail-type": "subscription-dontknow",
	"source": "prod.user.subscrition",
	"account": "4321",
	"time": "2019-12-22T18:43:48Z",
	"region": "us-east-1",
	"resources": [],
	"detail": {}
}`

const planCreatedJSON string = `{
	"version": "0",
	"id": "1234",
	"detail-type": "plan-created",
	"source": "prod.user.plan",
	"account": "4321",
	"time": "2019-12-22T18:43:48Z",
	"region": "us-east-1",
	"resources": [],
	"detail": {
		"uuid": "88bd0734-04cf-11ea-b52b-00155ddc81f9",
		"name": "Plan 1",
		"criteria": [{
			"uuid": "950470f4-04cf-11ea-8a5a-00155ddc81f8",
			"name": "Criteria 1",
			"max_count": 10,
			"model": "FREEMIUM",
			"action": "CREATE_CHALLENGE",
			"time_period": "WEEK"
		}]
	}
}`

const planUpdatedJSON string = `{
	"version": "0",
	"id": "1234",
	"detail-type": "plan-created",
	"source": "prod.user.plan",
	"account": "4321",
	"time": "2019-12-22T18:43:48Z",
	"region": "us-east-1",
	"resources": [],
	"detail": {
		"uuid": "88bd0734-04cf-11ea-b52b-00155ddc81f9",
		"name": "Plan 1",
		"criteria": [{
			"uuid": "950470f4-04cf-11ea-8a5a-00155ddc81f8",
			"name": "Criteria 1",
			"max_count": 10,
			"model": "FREEMIUM",
			"action": "CREATE_CHALLENGE",
			"time_period": "WEEK"
		},
		{
			"uuid": "950470f4-04cf-11ea-8a5a-00155ddc81f2",
			"name": "Criteria 2",
			"max_count": 10,
			"model": "FREEMIUM",
			"action": "CREATE_SHARE",
			"time_period": "WEEK"
		}]
	}
}`

const planInvalidJSON string = `{
	"version": "0",
	"id": "1234",
	"detail-type": "plan-dontknow",
	"source": "prod.user.subscrition",
	"account": "4321",
	"time": "2019-12-22T18:43:48Z",
	"region": "us-east-1",
	"resources": [],
	"detail": {}
}`

type EventBridgeSampleType int

const (
	CHALLENGE_CREATED_SAMPLE EventBridgeSampleType = iota
	CHALLENGE_SHARED_SAMPLE
	CHALLENGE_INVALID_SAMPLE
	SUBSCRIPTION_CREATED_SAMPLE
	SUBSCRIPTION_CANCELLED_SAMPLE
	SUBSCRIPTION_INVALID_SAMPLE
	PLAN_CREATED_SAMPLE
	PLAN_UPDATED_SAMPLE
)

func EventBridgeSample(sampleType EventBridgeSampleType) events.CloudWatchEvent {
	var sample string
	switch sampleType {
	case CHALLENGE_CREATED_SAMPLE:
		sample = challengeCreatedJSON
		break
	case CHALLENGE_SHARED_SAMPLE:
		sample = challengeSharedJSON
		break
	case CHALLENGE_INVALID_SAMPLE:
		sample = challengeInvalidJSON
		break
	case SUBSCRIPTION_CREATED_SAMPLE:
		sample = subscriptionCreatedJSON
		break
	case SUBSCRIPTION_CANCELLED_SAMPLE:
		sample = subscriptionCancelledJSON
		break
	case SUBSCRIPTION_INVALID_SAMPLE:
		sample = subscriptionInvalidJSON
		break
	case PLAN_CREATED_SAMPLE:
		sample = planCreatedJSON
		break
	case PLAN_UPDATED_SAMPLE:
		sample = planUpdatedJSON
		break
	default:
		log.Fatal("Sample type not found")
	}

	var event events.CloudWatchEvent
	err := json.Unmarshal([]byte(sample), &event)

	if err != nil {
		log.Fatal("Cannot mock event bridge event")
	}

	return event
}
