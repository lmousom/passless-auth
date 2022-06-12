package controllers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/lmousom/passless-auth/models/otpdata"
	"github.com/lmousom/passless-auth/utils"
)

var table = []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func SendOtp(phonenumber string) otpdata.SendOtpResponse {
	otp := GenerateOtp(6)
	ttl := 2 * 60 * 1000
	expiresIn := time.Now().UTC().UnixMilli() + int64(ttl)
	data := phonenumber + "." + otp + "." + strconv.FormatInt(expiresIn, 10)
	hash := utils.Encrypt([]byte(data))
	fullhash := hash + "." + strconv.FormatInt(expiresIn, 10)
	var response otpdata.SendOtpResponse = otpdata.SendOtpResponse{Status: "success", Message: "OTP sent successfully", Phone: phonenumber, Hash: fullhash}
	fmt.Println(otp)
	return response
}

func SendOtpHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var sendOtpRequest otpdata.SendOtpRequest
	err := decoder.Decode(&sendOtpRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")

	res := SendOtp(sendOtpRequest.Phone)
	json.NewEncoder(w).Encode(res)
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

func sendOtp(receiver string, otp string) {

	client := &http.Client{}
	var data = strings.NewReader(`Body=Hi there&From=+15017122661&To=+15558675310`)
	req, err := http.NewRequest("POST", "https://api.twilio.com/2010-04-01/Accounts/$TWILIO_ACCOUNT_SID/Messages.json", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("$TWILIO_ACCOUNT_SID", "$TWILIO_AUTH_TOKEN")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", bodyText)

}
