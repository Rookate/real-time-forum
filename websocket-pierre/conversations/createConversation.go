package conversations

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
	"websocket/server"
	"websocket/server/utils"
)

type Conversation struct {
	ConversationID string    `json:"conversation_uuid"`
	User1ID        string    `json:"user1_uuid"`
	User2ID        string    `json:"user2_uuid"`
	CreatedAt      time.Time `json:"created_at"`
}

func CreateConversation(db *sql.DB, r *http.Request, params map[string]interface{}) (*Conversation, error) {

	conversation_uuid, err := utils.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la génération du UUID: %v", err)
	}

	user1_uuid := "1"
	user2_uuid := "2"
	creationDate := time.Now()
	createConversationQuery := `INSERT INTO conversations (conversation_uuid, user1_uuid, user2_uuid, created_at) VALUES (?,?,?,?)`
	_, err = server.RunQuery(createConversationQuery, conversation_uuid, user1_uuid, user2_uuid, creationDate)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création de la conversation: %v", err)
	}

	newConversation := &Conversation{
		ConversationID: conversation_uuid,
		User1ID:        user1_uuid,
		User2ID:        user2_uuid,
		CreatedAt:      creationDate,
	}

	return newConversation, nil
}
