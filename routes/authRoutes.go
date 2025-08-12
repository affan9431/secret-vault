package routes

import (
	"github.com/affan9431/secret-vault/controllers"
	"github.com/gorilla/mux"
)

func AuthRoutes(router *mux.Router) {
	authRouter := router.PathPrefix("/api/user").Subrouter()
	authRouter.HandleFunc("/signUp", controllers.SignUpHandler).Methods("POST")
	authRouter.HandleFunc("/login", controllers.LoginHandler).Methods("POST")
}