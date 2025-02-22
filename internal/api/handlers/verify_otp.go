package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/lmousom/passless-auth/models/verifydata"
	"github.com/lmousom/passless-auth/utils"
)

var jwtKey = []byte("github.com/lmousom/passless-auth")

type Claims struct {
	Phone string `json:"phone"`
	jwt.RegisteredClaims
}

func VerifyOtp(verifyOtpRequest verifydata.VerifyOtpRequest, w http.ResponseWriter) verifydata.VerifyOtpResponse {
	f := func(c rune) bool {
		return c == '.'
	}
	extValue := strings.FieldsFunc(verifyOtpRequest.Hash, f)

	hashValue := extValue[0]
	expiresIn := extValue[1]

	now := time.Now()
	expiredInTime, err := utils.MsToTime(expiresIn)
	if err != nil {
		panic(err)
	}

	claims := &Claims{
		Phone: verifyOtpRequest.Phone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredInTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expiredInTime,
	})

	if now.After(expiredInTime) {
		return verifydata.VerifyOtpResponse{Status: "error", Message: "OTP Expired"}
	}

	if hashValue == utils.Encrypt([]byte(verifyOtpRequest.Phone+"."+verifyOtpRequest.Otp+"."+expiresIn)) {
		return verifydata.VerifyOtpResponse{Status: "success", Message: "OTP verified successfully"}
	} else {

		return verifydata.VerifyOtpResponse{Status: "error", Message: "OTP Invalid"}
	}

}

func VerifyOtpHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var verifyOtpRequest verifydata.VerifyOtpRequest
	err := decoder.Decode(&verifyOtpRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")

	res := VerifyOtp(verifyOtpRequest, w)

	json.NewEncoder(w).Encode(res)
}
