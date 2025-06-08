package twofa

type TwoFASettings struct {
	Phone     string `json:"phone"`
	Enabled   bool   `json:"enabled"`
	SecretKey string `json:"secret_key,omitempty"`
}

type Enable2FARequest struct {
	Phone string `json:"phone"`
}

type Enable2FAResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	SecretKey string `json:"secret_key"`
	QRCode    string `json:"qr_code"`
}

type Verify2FARequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type Verify2FAResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Disable2FARequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type Disable2FAResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
