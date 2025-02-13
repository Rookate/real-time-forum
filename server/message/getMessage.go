package message

import (
	"database/sql"
	"fmt"
	"forum/server"
	"net/http"
	"time"
)

func GetMessagesByConversations(db *sql.DB, r *http.Request, conversationUUID string) ([]Message, error) {
	getMessagesByConversationsQuery := `
	SELECT 
		m.message_uuid, 
		m.content, 
		m.created_at, 
		m.is_deleted, 
		m.sender_uuid,
		m.receiver_uuid,
		sender.username AS sender_username, 
		sender.profile_picture AS sender_profile_picture, 
		receiver.username AS receiver_username, 
		receiver.profile_picture AS receiver_profile_picture
		FROM messages m
		LEFT JOIN users sender ON m.sender_uuid = sender.user_uuid
		LEFT JOIN users receiver ON m.receiver_uuid = receiver.user_uuid
		WHERE m.conversation_uuid = ?
		ORDER BY m.created_at ASC;
	`

	rows, err := server.RunQuery(getMessagesByConversationsQuery, conversationUUID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des messages par conversation: %v", err)
	}

	var messages []Message
	for _, row := range rows {
		message := Message{
			ID:             row["message_uuid"].(string),
			ConversationID: conversationUUID,
			Content:        row["content"].(string),
			CreatedAt:      row["created_at"].(time.Time),
			IsDeleted:      row["is_deleted"].(bool),
			SenderID:       row["sender_uuid"].(string),
			ReceiverID:     row["receiver_uuid"].(string),
			// SenderUsername:   row["sender_username"].(string),
			// ReceiverUsername: row["receiver_username"].(string),
		}

		if data, ok := row["sender_profile_picture"]; ok && data != nil {
			message.SenderPicture = row["sender_profile_picture"].(string)
		}
		if data2, ok := row["receiver_profile_picture"]; ok && data2 != nil {
			message.ReceiverPicture = row["receiver_profile_picture"].(string)
		}
		if data3, ok := row["sender_username"]; ok && data3 != nil {
			message.SenderUsername = row["sender_username"].(string)
		}
		if data4, ok := row["receiver_username"]; ok && data4 != nil {
			message.ReceiverUsername = row["receiver_username"].(string)
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func GetSingleMessageByConversations(db *sql.DB, r *http.Request, conversationUUID string, messageUUID string) ([]Message, error) {
	getSingleMessageByConversationsQuery := `
	SELECT 
		m.message_uuid, 
		m.content, 
		m.created_at, 
		m.is_deleted, 
		m.sender_uuid,
		m.receiver_uuid,
		sender.username AS sender_username, 
		sender.profile_picture AS sender_profile_picture, 
		receiver.username AS receiver_username, 
		receiver.profile_picture AS receiver_profile_picture
	FROM messages m
	LEFT JOIN users sender ON m.sender_uuid = sender.user_uuid
	LEFT JOIN users receiver ON m.receiver_uuid = receiver.user_uuid
	WHERE m.conversation_uuid = ? AND m.message_uuid = ?
	LIMIT 1;
	`

	rows, err := server.RunQuery(getSingleMessageByConversationsQuery, conversationUUID, messageUUID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des messages par conversation: %v", err)
	}

	var singleMessage []Message
	for _, row := range rows {
		message := Message{
			ID:             row["message_uuid"].(string),
			ConversationID: conversationUUID,
			Content:        row["content"].(string),
			CreatedAt:      row["created_at"].(time.Time),
			IsDeleted:      row["is_deleted"].(bool),
			SenderID:       row["sender_uuid"].(string),
			ReceiverID:     row["receiver_uuid"].(string),
			// SenderUsername:   row["sender_username"].(string),
			// ReceiverUsername: row["receiver_username"].(string),
		}

		if data, ok := row["sender_profile_picture"]; ok && data != nil {
			message.SenderPicture = row["sender_profile_picture"].(string)
		}
		if data2, ok := row["receiver_profile_picture"]; ok && data2 != nil {
			message.ReceiverPicture = row["receiver_profile_picture"].(string)
		}
		if data3, ok := row["sender_username"]; ok && data3 != nil {
			message.SenderUsername = row["sender_username"].(string)
		}
		if data4, ok := row["receiver_username"]; ok && data4 != nil {
			message.ReceiverUsername = row["receiver_username"].(string)
		}

		singleMessage = append(singleMessage, message)

	}

	return singleMessage, nil
}
