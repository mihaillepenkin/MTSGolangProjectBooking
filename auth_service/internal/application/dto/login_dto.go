package dto

type LoginInputDto struct {
	Email   string `json:"email"`
	Passwrd string `json:"password"`
}

type LoginOutputDto struct {
	Msg    string `json:"msg"`
	Status string `json:"status"`
	Token  string `json:"token"`
}