package mailer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type Emailer interface {
	SendEmail(ctx context.Context, email *Email) error
}

type Emails struct {
	Emails []Email `json:"emails,omitempty"`
}

// taken from the BOP openapi spec
type Email struct {
	Subject    string   `json:"subject,omitempty"`
	Body       string   `json:"body,omitempty"`
	Recipients []string `json:"recipients,omitempty"`
	CcList     []string `json:"ccList,omitempty"`
	BccList    []string `json:"bccList,omitempty"`
	BodyType   string   `json:"bodyType,omitempty"`
}

func (e *Email) getBody() *types.Body {
	body := &types.Body{}

	if strings.ToLower(e.BodyType) == "html" {
		body.Html = &types.Content{Data: aws.String(e.Body)}
	} else {
		body.Text = &types.Content{Data: aws.String(e.Body)}
	}

	return body
}

func NewMailer() (Emailer, error) {
	var sender Emailer

	mod := os.Getenv("MAILER_MODULE")
	switch mod {
	case "aws":
		if cfg == nil {
			return nil, errors.New("config not initialized")
		}

		sender = &AwsSESEmailer{client: sesv2.NewFromConfig(*cfg)}
	default:
		return nil, fmt.Errorf("unsupported mailer module: %v", mod)
	}

	return sender, nil
}
