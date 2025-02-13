package apiConversations

import (
	"encoding/json"
	"forum/server"
	authentification "forum/server/api/login"
	"forum/server/conversations"
	"net/http"
)

func GetConversation(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user_uuid, _ := authentification.GetUserFromCookie(r)

	newConversation, err := conversations.GetConversations(server.Db, user_uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content", "application/json")
	json.NewEncoder(w).Encode(newConversation)
}
