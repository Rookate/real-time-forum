package apiMessages

import (
	"encoding/json"
	"forum/server"
	"forum/server/message"
	"net/http"
)

func CreateMessage(w http.ResponseWriter, r *http.Request) {
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

	var required []string = []string{"content", "conversation_uuid"}
	for _, key := range required {
		if _, ok := params[key]; !ok {
			http.Error(w, "Missing required field: "+key, http.StatusBadRequest)
			return
		}
	}

	newMessage, err := message.CreateMessage(server.Db, r, params)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newMessage)
}
