package handlers

import (
	"auth_service/internal/application/dto"
	"auth_service/internal/application/usecase"
	"encoding/json"
	"log"
	"net/http"
)

type HTTPHandler struct {
	userService *usecase.UserUseCase
}

func NewHandler(service *usecase.UserUseCase) *HTTPHandler {
	controller := new(HTTPHandler)
	controller.userService = service
	return controller
}

func (h *HTTPHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var input dto.CreateUserInputDto
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var output dto.CreateUserOutputDto
	output = h.userService.CreateUser(input.Name, input.Email, input.Passwrd, input.Role)
	json.NewEncoder(w).Encode(output)
}

func (h *HTTPHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var input dto.LoginInputDto
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var output dto.LoginOutputDto
	output = *h.userService.Login(input.Email, input.Passwrd)
	log.Printf("output: %v", output)
	json.NewEncoder(w).Encode(output)
}

func (h *HTTPHandler) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/register", h.RegisterHandler)
	mux.HandleFunc("/login", h.LoginHandler)

	return mux
}
