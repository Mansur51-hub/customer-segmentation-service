package handler

import (
	"encoding/json"
	"github.com/Mansur51-hub/customer-segmentation-service/service"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Handler struct {
	serv     *service.MyService
	validate *validator.Validate
}

func NewHandler(service *service.MyService) *Handler {
	return &Handler{serv: service, validate: validator.New()}
}

func NewResponse(w http.ResponseWriter, status int, msg any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(msg)
}
