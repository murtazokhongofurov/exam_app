package models

type UpdateAccessToken struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ResponseError struct {
	Error interface{} `json:"error"`
}

type ServerError struct {
	Status string `json:"status"`
	Message string `json:"message"`
}