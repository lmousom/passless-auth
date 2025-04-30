package sms

import (
	"fmt"
	"time"

	"github.com/lmousom/passless-auth/internal/config"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioService struct {
	client     *twilio.RestClient
	fromNumber string
	otpExpiry  time.Duration
}

func NewTwilioService(cfg *config.Config) (*TwilioService, error) {
	accountSID, err := cfg.GetDecryptedSMSAccountSID()
	if err != nil {
		return nil, fmt.Errorf("failed to get account SID: %w", err)
	}

	authToken, err := cfg.GetDecryptedSMSAuthToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get auth token: %w", err)
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})

	return &TwilioService{
		client:     client,
		fromNumber: cfg.SMS.FromNumber,
		otpExpiry:  cfg.Security.OTPExpiry,
	}, nil
}

func (s *TwilioService) SendOTP(phoneNumber, otp string) error {
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(phoneNumber)
	params.SetFrom(s.fromNumber)

	minutes := int(s.otpExpiry.Minutes())
	params.SetBody(fmt.Sprintf("Your OTP is: %s. Valid for %d minutes.", otp, minutes))

	_, err := s.client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	return nil
}
