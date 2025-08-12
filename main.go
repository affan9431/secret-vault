package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/affan9431/secret-vault/routes"
	"github.com/affan9431/secret-vault/storage"
	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "🔐 Secret Vault API")
}

func main() {
	var err error

	// Load .env file
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	if DB_PASSWORD == "" {
		panic("❌ Missing required environment variable: 'DB_PASSWORD'!")

	}
	// Connect to DB
	db, err = sql.Open("mysql", "root:"+DB_PASSWORD+"@tcp(127.0.0.1:3306)/securevault")
	if err != nil {
		panic("❌ Failed to connect to DB: " + err.Error())
	}

	storage.InitDB(db)

	// Close DB
	defer db.Close()

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	// Create new router
	router := mux.NewRouter()

	// Handle all requests
	router.HandleFunc("/", indexHandler)
	routes.AuthRoutes(router)
	routes.SecretRoutes(router)

	// Start server
	http.ListenAndServe(":"+PORT, router)
}

// In Go, only functions, variables, constants, or structs starting with an uppercase letter are exported (i.e., visible outside the package).
