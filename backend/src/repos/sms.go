package repos

import (
	"fmt"
	"os"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type SMS struct {
	client *twilio.RestClient
}

func NewSMS() *SMS {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: os.Getenv("TWILIO_ACCOUNT_SID"),
		Password: os.Getenv("TWILIO_AUTH_TOKEN"),
	})
	return &SMS{client: client}
}

func (s *SMS) SendOTP(phoneNumber, otpCode string) error {
	from := os.Getenv("TWILIO_PHONE_NUMBER")
	if from == "" {
		return fmt.Errorf("TWILIO_PHONE_NUMBER not configured")
	}

	message := fmt.Sprintf("Your verification code is: %s. Valid for 5 minutes.", otpCode)

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(phoneNumber)
	params.SetFrom(from)
	params.SetBody(message)

	_, err := s.client.Api.CreateMessage(params)
	return err
}
