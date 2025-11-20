package dto

type CreateUserInputDto struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Passwrd string `json:"password"`
	Role    string `json:"role"`
}

type CreateUserOutputDto struct {
	Msg    string `json:"msg"`
	Status   string `json:"status"`
}