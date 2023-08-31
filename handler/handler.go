package handler

import (
	"encoding/json"
	"github.com/Mansur51-hub/customer-segmentation-service/service"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Handler struct {
	serv     *service.Services
	validate *validator.Validate
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{serv: services, validate: validator.New()}
}

func NewResponse(w http.ResponseWriter, status int, msg any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(msg)
}
