package subscription

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/gusaul/go-dynamock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create Subscriptions Data Service Suite", func() {
	var (
		subscriptionRepository *SubscriptionRepository
		mock                   *dynamock.DynaMock
		db                     dynamodbiface.DynamoDBAPI
		subs                   map[string]*dynamodb.AttributeValue
	)

	JustBeforeEach(func() {
		var err error
		subscriptionRepository = NewSubscriptionRepository()
		db, mock = dynamock.New()
		subscriptionRepository.client = db
		subs, err = dynamodbattribute.MarshalMap(SubscriptionSample)
		if err != nil {
			panic(fmt.Sprintf("failed to DynamoDB marshal Record, %v", err))
		}
	})

	Context("Given subscription on DataService", func() {
		It("It should return the data on DynamoDB", func() {
			mock.ExpectScan().Table(subscriptionsTableName).WillReturns(dynamodb.ScanOutput{
				Items: []map[string]*dynamodb.AttributeValue{
					subs,
				},
			})

			subsResult, err := subscriptionRepository.GetByUserUUID(SubscriptionSample.UUID)
			Expect(err).To(BeNil())
			Expect(subsResult.UUID).To(Equal(SubscriptionSample.UUID))
			Expect(subsResult.Plan.Name).To((Equal(SubscriptionSample.Plan.Name)))
			Expect(subsResult.Provider).To((Equal(SubscriptionSample.Provider)))
			Expect(subsResult.Plan.Criteria[0].UUID).To(Equal(SubscriptionSample.Plan.Criteria[0].UUID))
		})
	})

	Context("Given subscription created on DataService", func() {
		It("It should create the data on DynamoDB", func() {
			mock.ExpectScan().Table(subscriptionsTableName).WillReturns(dynamodb.ScanOutput{
				Items: []map[string]*dynamodb.AttributeValue{
					subs,
				},
			})

			mock.ExpectTransactWriteItems().WillReturns(dynamodb.TransactWriteItemsOutput{})

			err := subscriptionRepository.Store(SubscriptionSample)
			Expect(err).To(BeNil())
		})
	})

	Context("Given subscription created on DataService", func() {
		It("It should create the data on DynamoDB if no prev sub was found", func() {
			mock.ExpectScan().Table(subscriptionsTableName).WillReturns(dynamodb.ScanOutput{
				Items: []map[string]*dynamodb.AttributeValue{},
			})

			mock.ExpectPutItem().WillReturns(dynamodb.PutItemOutput{})

			err := subscriptionRepository.Store(SubscriptionSample)
			Expect(err).To(BeNil())
		})
	})

	Context("Given subscription cancelled on DataService", func() {
		It("It should cancel the data on DynamoDB", func() {
			mock.ExpectUpdateItem().ToTable(subscriptionsTableName).WillReturns(dynamodb.UpdateItemOutput{
				Attributes: subs,
			})

			err := subscriptionRepository.CancelSubscription(SubscriptionSample.UUID)
			Expect(err).To(BeNil())
		})
	})
})

func TestSubscriptionsDataService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Create Subscriptions Data Service Suite")
}
