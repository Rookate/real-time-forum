package message

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
	"websocket/server"
)

func GetMessagesByConversations(db *sql.DB, r *http.Request, conversationUUID string) ([]Message, error) {

	getMessagesByConversationsQuery := `
	SELECT 
    m.message_uuid, 
    m.message, 
    m.created_at, 
    m.is_deleted, 
	m.sender_uuid,
    m.receiver_uuid,
    u1.username AS sender_username, 
    u1.profil_picture AS sender_profile_picture, 
    u2.username AS receiver_username, 
    u2.profil_picture AS receiver_profile_picture
	FROM messages m
	JOIN users u1 ON m.sender_uuid = u1.user_uuid
	JOIN users u2 ON m.receiver_uuid = u2.user_uuid
	WHERE m.conversation_uuid = ?
	ORDER BY m.created_at ASC`

	// conversationUUID := params["conversation_uuid"].(string)

	param := conversationUUID

	rows, err := server.RunQuery(getMessagesByConversationsQuery, param)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des messages par conversations: %v", err)
	}

	var messages []Message
	for _, row := range rows {
		message := Message{
			ID:             row["message_uuid"].(string),
			ConversationID: conversationUUID,
			Message:        row["message"].(string),
			CreatedAt:      row["created_at"].(time.Time),
			IsDeleted:      row["is_deleted"].(bool),
			SenderID:       row["sender_uuid"].(string),
			ReceiverID:     row["receiver_uuid"].(string),
		}
		messages = append(messages, message)
	}
	return messages, nil

}
