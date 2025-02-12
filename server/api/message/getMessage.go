package apiMessages

import (
	"encoding/json"
	"forum/server"
	"forum/server/message"
	"net/http"
)

func GetMessageByConversation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var params map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	conversation_uuid, ok := params["conversation_uuid"].(string)
	if !ok || conversation_uuid == "" {
		http.Error(w, "Missing or invalid conversation_uuid", http.StatusBadRequest)
		return
	}
	conversationUUIDdata, err := message.GetMessagesByConversations(server.Db, r, conversation_uuid)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content", "application/json")
	if err := json.NewEncoder(w).Encode(conversationUUIDdata); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
