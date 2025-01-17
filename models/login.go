package models

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Status       string `json:"status"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
