package message

import (
	"database/sql"
	"fmt"
	"forum/server"
	authentification "forum/server/api/login"
	posts "forum/server/utils"
	"net/http"
	"time"
)

type Message struct {
	ID               string    `json:"message_uuid"`
	ConversationID   string    `json:"conversation_uuid"`
	SenderID         string    `json:"sender_uuid"`
	ReceiverID       string    `json:"receiver_uuid"`
	Content          string    `json:"content"`
	CreatedAt        time.Time `json:"created_at"`
	SenderPicture    string    `json:"sender_profile_picture"`
	SenderUsername   string    `json:"sender_username"`
	ReceiverUsername string    `json:"receiver_username"`
	ReceiverPicture  string    `json:"receiver_profile_picture"`
	IsDeleted        bool      `json:"is_deleted"`
}

func CreateMessage(db *sql.DB, r *http.Request, params map[string]interface{}) (*Message, error) {

	message_uuid, err := posts.GenerateUUID()

	if err != nil {
		return nil, fmt.Errorf("erreur lors de la génération du UUID: %v", err)
	}

	conversation_uuid, ok := params["conversation_uuid"].(string)
	if !ok || conversation_uuid == "" {
		return nil, fmt.Errorf("UUID de conversation manquant")
	}

	user2, _ := params["receiver_uuid"].(string)
	sender_uuid, _ := authentification.GetUserFromCookie(r)
	receiver_uuid := user2

	content, _ := params["content"].(string)
	creationDate := time.Now()
	is_deleted := false

	createMessageQuery := `INSERT INTO messages (message_uuid, conversation_uuid, sender_uuid, receiver_uuid, content, created_at, is_deleted) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = server.RunQuery(createMessageQuery, message_uuid, conversation_uuid, sender_uuid, receiver_uuid, content, creationDate, is_deleted)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création du post: %v", err)
	}

	newMessage := &Message{
		ID:             message_uuid,
		ConversationID: conversation_uuid,
		SenderID:       sender_uuid,
		ReceiverID:     receiver_uuid,
		Content:        content,
		CreatedAt:      creationDate,
		IsDeleted:      is_deleted,
	}

	return newMessage, nil
}
