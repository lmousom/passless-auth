package verifydata

type VerifyOtpRequest struct {
	Phone string `json:"phone"`
	Hash  string `json:"hash"`
	Otp   string `json:"otp"`
}

type VerifyOtpResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
