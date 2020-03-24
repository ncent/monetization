package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type ExpressionScanner struct {
	client dynamodbiface.DynamoDBAPI
}

func NewExpressionScanner(dynamodbClient dynamodbiface.DynamoDBAPI) *ExpressionScanner {
	return &ExpressionScanner{client: dynamodbClient}
}

func (es ExpressionScanner) ScanWithExpression(tableName string, conditions ...expression.ConditionBuilder) (*dynamodb.ScanOutput, error) {
	exprBuilder := expression.NewBuilder()

	for _, condition := range conditions {
		exprBuilder = exprBuilder.WithFilter(condition)
	}

	expr, err := exprBuilder.Build()

	if err != nil {
		return nil, err
	}

	queryInput := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	return es.client.Scan(queryInput)
}
