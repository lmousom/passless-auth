package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lmousom/passless-auth/controllers"
)

func main() {

	r := AuthRouter()
	http.ListenAndServe(":8080", r)

}

func AuthRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/sendOtp", controllers.SendOtpHandler).Methods("POST")
	r.HandleFunc("/api/v1/verifyOtp", controllers.VerifyOtpHandler).Methods("POST")
	r.HandleFunc("/api/v1/login", controllers.VerificationHandler).Methods("GET")
	r.HandleFunc("/api/v1/refreshToken", controllers.RefreshTokenHandler).Methods("POST")

	return r
}
