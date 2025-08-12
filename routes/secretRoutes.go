package routes

import (
	"net/http"

	"github.com/affan9431/secret-vault/controllers"
	"github.com/affan9431/secret-vault/middleware"
	"github.com/gorilla/mux"
)

func SecretRoutes(router *mux.Router) {
	secretRouter := router.PathPrefix("/api/secrets").Subrouter()
	secretRouter.Handle("/create-secret", middleware.AuthMiddleware(http.HandlerFunc(controllers.CreateSecretHandler)))
	secretRouter.Handle("/get-secret", middleware.AuthMiddleware(http.HandlerFunc(controllers.GetSecretHandler)))
}
