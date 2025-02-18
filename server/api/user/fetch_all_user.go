package user

import (
	"encoding/json"
	"forum/server"
	"log"
	"net/http"
)

func FetchAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer tous les posts
	userData, err := FetchAllUsers(server.Db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(userData)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
