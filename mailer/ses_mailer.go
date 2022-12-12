package mailer

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

var cfg *aws.Config

func InitConfig() error {
	mod := os.Getenv("MAILER_MODULE")
	switch mod {
	case "aws":
		config, err := config.LoadDefaultConfig(context.Background(),
			config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
				Value: aws.Credentials{
					AccessKeyID:     os.Getenv("AWS_ACCESS_KEY"),
					SecretAccessKey: os.Getenv("AWS_SECRET_KEY"),
					Source:          "fedrampbop",
				},
			}))
		if err != nil {
			return err
		}

		cfg = &config
	default:
		return fmt.Errorf("unsupported mailer module: %v", mod)
	}

	return nil
}

var _ = (Emailer)(&AwsSESEmailer{})

type AwsSESEmailer struct {
	client *sesv2.Client
}

func (s *AwsSESEmailer) SendEmail(ctx context.Context, email *Email) error {
	out, err := s.client.SendEmail(ctx, &sesv2.SendEmailInput{
		// what is this? will need to be validated in AWS
		FromEmailAddress: aws.String("no-reply@redhat.com"),
		Destination: &types.Destination{
			// TODO: integrate username lookups? the docs indicate that but not
			// sure if it would actually be necessary here.
			// TODO: support for "\"Real Name\" user@example.com" sending, right
			// now AWS wants _just_ the email so we will have to sanitize the input
			ToAddresses:  email.Recipients,
			CcAddresses:  email.CcList,
			BccAddresses: email.BccList,
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{Data: aws.String(email.Subject)},
				Body:    email.getBody(),
			}},
	})
	if err != nil {
		return err
	}

	log.Printf("Sent message %#v successfully, msg id: %v", email, out.MessageId)
	return nil
}
