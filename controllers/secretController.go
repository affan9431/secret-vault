package controllers

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
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

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "✅ Successfully created secret!",
	})
}

// TODO: Implement this function
func GetSecretHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := storage.DB.Query("SELECT id, title, secret, tags, extra_data FROM user_secrets WHERE user_id = ?", r.URL.Query().Get("id"))

	if err != nil {
		fmt.Println(err)
	}

	keyHex := os.Getenv("VAULT_ENC_KEY")
	key, _ := hex.DecodeString(keyHex)

	type SecretResponse struct {
		Id        int64  `json:"id"`
		Title     string `json:"title"`
		Secret    string `json:"secret"`
		Tags      string `json:"tags"`
		ExtraData string `json:"extra_data,omitempty"`
	}

	var secrets []SecretResponse

	for rows.Next() {
		var title, secret, tags, extra_data []byte
		var id int64

		if err := rows.Scan(&id, &title, &secret, &tags, &extra_data); err != nil {
			fmt.Println(err)
		}

		decryptedTitle, _ := utils.Decrypt([]byte(title), key)
		decryptedSecret, _ := utils.Decrypt([]byte(secret), key)
		decryptedTag, _ := utils.Decrypt([]byte(tags), key)

		var decryptedExtraData []byte
		if len(extra_data) > 0 {
			decryptedExtraData, _ = utils.Decrypt(extra_data, key)
		}

		secrets = append(secrets, SecretResponse{
			Id:        id,
			Title:     string(decryptedTitle),
			Secret:    string(decryptedSecret),
			Tags:      string(decryptedTag),
			ExtraData: string(decryptedExtraData),
		})

		fmt.Println("ID: " + fmt.Sprint(id))
		fmt.Println("Title: " + string(decryptedTitle))
		fmt.Println("Secret: " + string(decryptedSecret))
		fmt.Println("Tags: " + string(decryptedTag))
		fmt.Println("Extra Data: " + string(decryptedExtraData))
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"secrets": secrets,
	})

}

func UpdateSecretHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	type UpdateSecret struct {
		Title     string  `json:"title"`
		Secret    string  `json:"secret"`
		Tags      string  `json:"tags"`
		ExtraData *string `json:"extra_data,omitempty"`
	}

	var updateSecret UpdateSecret
	json.NewDecoder(r.Body).Decode(&updateSecret)

	extra_data := updateSecret.ExtraData

	keyHex := os.Getenv("VAULT_ENC_KEY")
	key, _ := hex.DecodeString(keyHex)

	encrpytedTitle, err := utils.Encrypt([]byte(updateSecret.Title), key)

	if err != nil {
		fmt.Println(err)
	}
	encrpytedSecret, err := utils.Encrypt([]byte(updateSecret.Secret), key)

	if err != nil {
		fmt.Println(err)
	}
	encrpytedTag, err := utils.Encrypt([]byte(updateSecret.Tags), key)

	if err != nil {
		fmt.Println(err)
	}
	encrpytedExtraData, err := utils.Encrypt([]byte(*updateSecret.ExtraData), key)

	if err != nil {
		fmt.Println(err)
	}

	if extra_data != nil && *extra_data != "" {

		_, err := storage.DB.Exec("UPDATE user_secrets SET title=?, secret=?, tags=?, extra_data=? WHERE id=?", encrpytedTitle, encrpytedSecret, encrpytedTag, encrpytedExtraData, id)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		_, err := storage.DB.Exec("UPDATE user_secrets SET title=?, secret=?, tags=? WHERE id=?", encrpytedTitle, encrpytedSecret, encrpytedTag, id)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "✅ Successfully update secret!",
	})

}

func DeleteSecretHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	result, err := storage.DB.Exec("DELETE FROM user_secrets WHERE id=?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	count, _ := result.RowsAffected()
	if count == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "❌ No secret found with that title for this user",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "✅ Deleted successfully!",
	})
}
