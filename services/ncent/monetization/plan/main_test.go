package plan

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/gusaul/go-dynamock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create Plans Data Service Suite", func() {
	var (
		planRepository *PlanRepository
		mock           *dynamock.DynaMock
	)

	JustBeforeEach(func() {
		var db dynamodbiface.DynamoDBAPI
		planRepository = NewPlanRepository()
		db, mock = dynamock.New()
		planRepository.client = db
	})

	Context("Given subscription on DataService", func() {
		It("It should return the data on DynamoDB", func() {
			uuid := "55109e64-1fa1-4749-ab4f-90cb781532a7"

			mock.ExpectQuery().Table(planTableName).WillReturns(dynamodb.QueryOutput{
				Items: []map[string]*dynamodb.AttributeValue{
					{
						"uuid": {
							S: aws.String("55109e64-1fa1-4749-ab4f-90cb781532a7"),
						},
						"name": {
							S: aws.String("Plan1"),
						},
					},
				},
			})

			p, err := planRepository.GetByUUID(uuid)
			Expect(err).To(BeNil())
			Expect(p.UUID).To(Equal(uuid))
			Expect(p.Name).To(Equal("Plan1"))
		})
	})

	Context("Given plan created on DataService", func() {
		It("It should create the data on DynamoDB", func() {
			mock.ExpectPutItem().ToTable(planTableName).WillReturns(dynamodb.PutItemOutput{})
			err := planRepository.Store(PlanSample)
			Expect(err).To(BeNil())
		})
	})

	Context("Given exiting plan on DataService", func() {
		It("It should be updated on DynamoDB", func() {
			uuid := "e55c8798-5493-4e69-8d58-74de4f1ed138"
			mock.ExpectUpdateItem().ToTable(planTableName).WillReturns(dynamodb.UpdateItemOutput{})
			err := planRepository.Update(uuid, "set name = name", PlanSample)
			Expect(err).To(BeNil())
		})
	})
})

func TestPlansDataService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Create Plans Data Service Suite")
}
