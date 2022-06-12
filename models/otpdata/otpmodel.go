package otpdata

type SendOtpResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Phone   string `json:"phone"`
	Hash    string `json:"hash"`
}

type SendOtpRequest struct {
	Phone string `json:"phone"`
}
