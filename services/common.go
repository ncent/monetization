package services

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/private/protocol/jsonrpc"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ses"
)

type ISQLDB interface {
	Begin() (*sql.Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Close() error
	Conn(ctx context.Context) (*sql.Conn, error)
	Driver() driver.Driver
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Ping() error
	PingContext(ctx context.Context) error
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	SetConnMaxLifetime(d time.Duration)
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
	Stats() sql.DBStats
}

type ICloudWatchEvents interface {
	PutEvents(putEventsInput *cloudwatchevents.PutEventsInput) (*cloudwatchevents.PutEventsOutput, error)
}

type ISESClient interface {
	SendEmail(input *ses.SendEmailInput) (*ses.SendEmailOutput, error)
	ListIdentities(input *ses.ListIdentitiesInput) (*ses.ListIdentitiesOutput, error)
	VerifyEmailAddress(input *ses.VerifyEmailAddressInput) (*ses.VerifyEmailAddressOutput, error)
}

const TxAwareErrorUnmarshallerName = "awssdk.jsonrpc.TxAwareErrorUnmarshaller"

// New creates a new instance of the DynamoDB client with a session.
// The client's behaviour is same as what is returned by dynamodb.New(), except for richer error reasons.
func NewTxErrorAwareDynamoDBClient(p client.ConfigProvider, cfgs ...*aws.Config) *dynamodb.DynamoDB {
	c := dynamodb.New(p, cfgs...)
	// NOTE: Ignore if swap failed. Returning nil might fail app startup which is worse than inadequate error details.
	c.Handlers.UnmarshalError.Swap(jsonrpc.UnmarshalErrorHandler.Name, request.NamedHandler{
		Name: TxAwareErrorUnmarshallerName,
		Fn:   TxAwareUnmarshalError,
	})
	return c
}

// TxAwareUnmarshalError unmarshals an error response for a JSON RPC service.
// This is exactly same as jsonrpc.UnmarshalError, except for attempt to parse CancellationReasons
func TxAwareUnmarshalError(req *request.Request) {
	defer req.HTTPResponse.Body.Close()

	var jsonErr jsonTxErrorResponse
	err := json.NewDecoder(req.HTTPResponse.Body).Decode(&jsonErr)
	if err == io.EOF {
		req.Error = awserr.NewRequestFailure(
			awserr.New(request.ErrCodeSerialization, req.HTTPResponse.Status, nil),
			req.HTTPResponse.StatusCode,
			req.RequestID,
		)
		return
	} else if err != nil {
		req.Error = awserr.NewRequestFailure(
			awserr.New(request.ErrCodeSerialization,
				"failed decoding JSON RPC error response", err),
			req.HTTPResponse.StatusCode,
			req.RequestID,
		)
		return
	}

	codes := strings.SplitN(jsonErr.Code, "#", 2)
	req.Error = newTxRequestError(
		awserr.New(codes[len(codes)-1], jsonErr.Message, nil),
		req.HTTPResponse.StatusCode,
		req.RequestID,
		jsonErr.CancellationReasons,
	)
}

// A RequestFailure is an interface to extract request failure information from an Error.
type TxRequestFailure interface {
	awserr.RequestFailure
	CancellationReasons() []dynamodb.CancellationReason
}

type jsonTxErrorResponse struct {
	Code                string                        `json:"__type"`
	Message             string                        `json:"message"`
	CancellationReasons []dynamodb.CancellationReason `json:"CancellationReasons"`
}

// So that the Error interface type can be included as an anonymous field
// in the requestError struct and not conflict with the error.Error() method.
type awsError awserr.Error

// A TxRequestError wraps a request or service error.
// TxRequestError is awserr.requestError with additional cancellationReasons field
type TxRequestError struct {
	awsError
	statusCode          int
	requestID           string
	cancellationReasons []dynamodb.CancellationReason
}

func newTxRequestError(err awserr.Error, statusCode int, requestID string, cancellationReasons []dynamodb.CancellationReason) TxRequestFailure {
	return &TxRequestError{
		awsError:            err,
		statusCode:          statusCode,
		requestID:           requestID,
		cancellationReasons: cancellationReasons,
	}
}

func (r TxRequestError) Error() string {
	extra := fmt.Sprintf("status code: %d, request id: %s",
		r.statusCode, r.requestID)
	return awserr.SprintError(r.Code(), r.Message(), extra, r.OrigErr())
}

func (r TxRequestError) String() string {
	return r.Error()
}

func (r TxRequestError) StatusCode() int {
	return r.statusCode
}

func (r TxRequestError) RequestID() string {
	return r.requestID
}

func (r TxRequestError) CancellationReasons() []dynamodb.CancellationReason {
	return r.cancellationReasons
}
