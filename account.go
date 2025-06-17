package ember_backend_api_gateway

type SendOtpRequest struct {
	Email string `json:"email"`
}

type VerifyOtpRequest struct {
	Otp string `json:"otp"`
}

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}
