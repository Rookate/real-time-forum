package message

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
	"websocket/server"
	"websocket/server/utils"
)

type Message struct {
	ID             string    `json:"message_uuid"`
	ConversationID string    `json:"conversation_uuid"`
	SenderID       string    `json:"sender_uuid"`
	ReceiverID     string    `json:"receiver_uuid"`
	Message        string    `json:"message"`
	CreatedAt      time.Time `json:"created_at"`
	IsDeleted      bool      `json:"is_deleted"`
}

func CreateMessage(db *sql.DB, r *http.Request, params map[string]interface{}) (*Message, error) {

	message_uuid, err := utils.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la génération du UUID: %v", err)
	}

	conversation_uuid, ok := params["conversation_uuid"].(string)
	if !ok || conversation_uuid == "" {
		return nil, fmt.Errorf("UUID de conversation manquant")
	}

	fmt.Println("UUID Récupéré lors de la création du message:", conversation_uuid)

	sender_uuid := "1"
	receiver_uuid := "2"

	message, _ := params["content"].(string)
	fmt.Println("Message 1", message)
	creationDate := time.Now()
	is_deleted := false

	createMessageQuery := `INSERT INTO messages (message_uuid, conversation_uuid, sender_uuid, receiver_uuid, message, created_at, is_deleted) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = server.RunQuery(createMessageQuery, message_uuid, conversation_uuid, sender_uuid, receiver_uuid, message, creationDate, is_deleted)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création du post: %v", err)
	}

	newMessage := &Message{
		ID:             message_uuid,
		ConversationID: conversation_uuid,
		SenderID:       sender_uuid,
		ReceiverID:     receiver_uuid,
		Message:        message,
		CreatedAt:      creationDate,
		IsDeleted:      is_deleted,
	}

	return newMessage, nil
}
