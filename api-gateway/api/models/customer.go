package models

type CustomerRegister struct {
	FullName  string       `json:"full_name"`
	Bio       string       `json:"bio"`
	Email     string       `json:"email"`
	Password  string       `json:"password"`
	Addresses []AddressReq `json:"addresses"`
}

type AddressReq struct {
	Country string `json:"country"`
	Street  string `json:"street"`
}

type Customer struct {
	FullName  string    `json:"full_name"`
	Bio       string    `json:"bio"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Addresses []Address `json:"addresses"`
}
type Address struct {
	Id      string `json:"id"`
	OwnerId string `json:"owner_id"`
	Country string `json:"country"`
	Street  string `json:"street"`
}

type VerifyResponse struct {
	Id           string    `json:"id"`
	FullName     string    `json:"full_name"`
	Bio          string    `json:"bio"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	AccessToken  string    `json:"accsee_token"`
	RefreshToken string    `json:"refresh_token"`
	Addresses    []Address `json:"addresses"`
}

type CustomerUpdateReq struct {
	Id       string `json:"id"`
	FullName string `json:"full_name"`
	Bio      string `json:"bio"`
	Email    string `json:"email"`
}

type CustomerLogin struct {
	FullName     string    `json:"full_name"`
	Bio          string    `json:"bio"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	RefreshToken string    `json:"refresh_token"`
	AccessToken  string    `json:"access_token"`
	Addresses    []Address `json:"addresses"`
}
