package apiConversations

import (
	"encoding/json"
	"net/http"
	"websocket/server"
	"websocket/server/conversations"
)

func CreateConversation(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var params map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	newConverstion, err := conversations.CreateConversation(server.DB, r, params)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content", "application/json")
	json.NewEncoder(w).Encode(newConverstion)
}
