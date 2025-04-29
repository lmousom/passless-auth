package handlers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/lmousom/passless-auth/internal/errors"
	"github.com/lmousom/passless-auth/internal/middleware"
	"github.com/lmousom/passless-auth/models/otpdata"
	"github.com/lmousom/passless-auth/utils"
)

var table = []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func SendOtp(phonenumber string) (*otpdata.SendOtpResponse, error) {
	if phonenumber == "" {
		return nil, errors.NewInvalidRequest("Phone number is required", nil)
	}

	otp := GenerateOtp(6)
	ttl := 2 * 60 * 1000
	expiresIn := time.Now().UTC().UnixMilli() + int64(ttl)
	data := phonenumber + "." + otp + "." + strconv.FormatInt(expiresIn, 10)
	hash := utils.Encrypt([]byte(data))
	fullhash := hash + "." + strconv.FormatInt(expiresIn, 10)

	response := &otpdata.SendOtpResponse{
		Status:  "success",
		Message: "OTP sent successfully",
		Phone:   phonenumber,
		Hash:    fullhash,
	}

	// TODO: Implement actual SMS sending
	fmt.Println(otp)
	return response, nil
}

func SendOtpHandler(w http.ResponseWriter, r *http.Request) {
	var sendOtpRequest otpdata.SendOtpRequest
	if err := json.NewDecoder(r.Body).Decode(&sendOtpRequest); err != nil {
		middleware.ErrorResponse(w, errors.NewInvalidRequest("Invalid request body", err))
		return
	}

	response, err := SendOtp(sendOtpRequest.Phone)
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
