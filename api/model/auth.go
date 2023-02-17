package model

type Token struct {
	Value string
}

type GenerateTokenRequest struct {
	Email    string
	Password string
}
