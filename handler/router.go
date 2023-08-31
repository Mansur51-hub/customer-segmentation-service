package handler

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Router struct {
	Router *mux.Router
}

func NewRouter(h *Handler) *Router {
	router := mux.NewRouter()

	router.HandleFunc("/users", h.CreateUserSegments).Methods("POST")
	router.HandleFunc("/users", h.GetUserSegments).Methods("GET")
	router.HandleFunc("/segments", h.CreateSegment).Methods("POST")
	router.HandleFunc("/segments", h.DeleteSegment).Methods("DELETE")
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	router.HandleFunc("/operations", h.GetOperations).Methods("GET")

	return &Router{Router: router}
}
