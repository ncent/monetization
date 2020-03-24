package client

import (
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/sirupsen/logrus"
	"gitlab.com/ncent/monetization/services"
)

var log *logrus.Entry

func init() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	log = logger.WithFields(logrus.Fields{
		"service": "service.aws.ses.client",
	})
}

type ISESService interface {
	SendEmail(er EmailRequest) error
}

type SESService struct {
	client services.ISESClient
}

func NewSESService() *SESService {
	cfg := &aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}

	return &SESService{
		client: ses.New(session.Must(session.NewSession(cfg))),
	}
}

func (sess SESService) SendEmail(er EmailRequest) error {
	log.Infof("er: %+v", er)

	err := sess.validateEmails(er.Sender)
	if err != nil {
		return sess.checkMailerError(err)
	}

	input := sess.createEmailInput(er)
	log.Infof("input: %+v", input)

	result, err := sess.client.SendEmail(input)
	if err != nil {
		return sess.checkMailerError(err)
	}

	log.Infof("Email Sent to address: %v", er.Recipient)
	log.Infof("Result: %+v", result)

	return nil
}

func (sess SESService) validateEmails(emails ...string) error {
	identitiesResult, _ := sess.client.ListIdentities(
		&ses.ListIdentitiesInput{
			IdentityType: aws.String("EmailAddress"),
		},
	)

	log.Infof("Identities found: %+v", identitiesResult.Identities)
	identities := make(map[string]string, len(identitiesResult.Identities))

	for _, s := range identitiesResult.Identities {
		identities[*s] = strings.ToLower(*s)
	}

	log.Infof("Identities map: %+v", identities)
	log.Infof("Comparing with emails: %+v", emails)

	for _, email := range emails {
		if _, exists := identities[strings.ToLower(email)]; !exists {
			_, err := sess.client.VerifyEmailAddress(&ses.VerifyEmailAddressInput{EmailAddress: aws.String(email)})
			if err != nil {
				return sess.checkMailerError(err)
			}
		}
	}

	return nil
}

func (sess SESService) createEmailInput(er EmailRequest) *ses.SendEmailInput {
	return &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(er.Recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(er.Html),
				},
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(er.Body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(er.Subject),
			},
		},
		Source: aws.String(er.Sender),
	}
}

func (sess SESService) checkMailerError(err error) error {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case ses.ErrCodeMessageRejected:
			log.Errorln(ses.ErrCodeMessageRejected, aerr.Error())
		case ses.ErrCodeMailFromDomainNotVerifiedException:
			log.Errorln(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
		case ses.ErrCodeConfigurationSetDoesNotExistException:
			log.Errorln(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
		default:
			log.Errorln(aerr.Error())
		}
	} else {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		log.Errorln(err.Error())
	}

	return err
}
