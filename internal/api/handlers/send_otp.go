package handlers

import (
	"crypto/rand"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/lmousom/passless-auth/internal/errors"
	"github.com/lmousom/passless-auth/internal/middleware"
	"github.com/lmousom/passless-auth/internal/services/sms"
	"github.com/lmousom/passless-auth/models/otpdata"
	"github.com/lmousom/passless-auth/utils"
)

var table = []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

type SendOtpHandler struct {
	smsService *sms.TwilioService
}

func NewSendOtpHandler(smsService *sms.TwilioService) *SendOtpHandler {
	return &SendOtpHandler{
		smsService: smsService,
	}
}

func (h *SendOtpHandler) SendOtp(phonenumber string) (*otpdata.SendOtpResponse, error) {
	if phonenumber == "" {
		return nil, errors.NewInvalidRequest("Phone number is required", nil)
	}

	otp := GenerateOtp(6)
	ttl := 2 * 60 * 1000
	expiresIn := time.Now().UTC().UnixMilli() + int64(ttl)
	data := phonenumber + "." + otp + "." + strconv.FormatInt(expiresIn, 10)
	hash := utils.Encrypt([]byte(data))
	fullhash := hash + "." + strconv.FormatInt(expiresIn, 10)

	// Send OTP via Twilio
	if err := h.smsService.SendOTP(phonenumber, otp); err != nil {
		return nil, errors.NewInternalServer("Failed to send OTP", err)
	}

	response := &otpdata.SendOtpResponse{
		Status:  "success",
		Message: "OTP sent successfully",
		Phone:   phonenumber,
		Hash:    fullhash,
	}

	return response, nil
}

func (h *SendOtpHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var sendOtpRequest otpdata.SendOtpRequest
	if err := json.NewDecoder(r.Body).Decode(&sendOtpRequest); err != nil {
		middleware.ErrorResponse(w, errors.NewInvalidRequest("Invalid request body", err))
		return
	}

	response, err := h.SendOtp(sendOtpRequest.Phone)
	if err != nil {
		middleware.ErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		middleware.ErrorResponse(w, errors.NewInternalServer("Failed to encode response", err))
		return
	}
}

func GenerateOtp(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}
