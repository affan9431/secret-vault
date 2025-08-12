package controllers

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/affan9431/secret-vault/models"
	"github.com/affan9431/secret-vault/storage"
	"github.com/affan9431/secret-vault/utils"
)

func CreateSecretHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("id")
	var secret models.Secrets
	json.NewDecoder(r.Body).Decode(&secret)
	extra_data := secret.ExtraData

	keyHex := os.Getenv("VAULT_ENC_KEY")
	key, _ := hex.DecodeString(keyHex)

	encrpytedTitle, err := utils.Encrypt([]byte(secret.Title), key)

	if err != nil {
		fmt.Println(err)
	}
	encrpytedSecret, err := utils.Encrypt([]byte(secret.Secret), key)

	if err != nil {
		fmt.Println(err)
	}
	encrpytedTag, err := utils.Encrypt([]byte(secret.Tags), key)

	if err != nil {
		fmt.Println(err)
	}
	encrpytedExtraData, err := utils.Encrypt([]byte(*secret.ExtraData), key)

	if err != nil {
		fmt.Println(err)
	}

	if extra_data != nil && *extra_data != "" {
		_, err := storage.DB.Exec("INSERT INTO user_secrets (user_id, title, secret, tags, extra_data) VALUES (?, ?, ?, ?, ?)", query, encrpytedTitle, encrpytedSecret, encrpytedTag, encrpytedExtraData)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else {
		_, err := storage.DB.Exec("INSERT INTO user_secrets (user_id, title, secret, tags) VALUES (?, ?, ?, ?)", query, encrpytedTitle, encrpytedSecret, encrpytedTag)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

	fmt.Fprintf(w, "Successfully created secret!")
}

// TODO: Implement this function
func GetSecretHandler(w http.ResponseWriter, r *http.Request) {

}
