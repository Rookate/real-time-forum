package conversations

import (
	"database/sql"
	"fmt"
	authentification "forum/server/api/login"
	posts "forum/server/utils"
	"net/http"
	"time"
)

type Conversation struct {
	ConversationID   string    `json:"conversation_uuid"`
	User1ID          string    `json:"sender"`
	User2ID          string    `json:"receiver"`
	ReceiverUsername string    `json:"receiver_username"`
	ReceiverPicture  string    `json:"receiver_profile_picture"`
	CreatedAt        time.Time `json:"created_at"`
}

func CreateConversation(db *sql.DB, r *http.Request, params map[string]interface{}) (*Conversation, error) {
	user2, _ := params["user_uuid"].(string)
	currentUser, _ := authentification.GetUserFromCookie(r)

	var otherUserPicture, otherUserUsername, conversationID string

	// Vérifier si une conversation existe déjà
	existingConversation := `
    SELECT 
        c.conversation_uuid,
        CASE 
            WHEN c.sender = ? THEN u2.profile_picture
            ELSE u1.profile_picture
        END as other_picture,
        CASE 
            WHEN c.sender = ? THEN u2.username
            ELSE u1.username
        END as other_username
    FROM 
        conversations c
    JOIN 
        users u1 ON u1.user_uuid = c.sender
    JOIN 
        users u2 ON u2.user_uuid = c.reciever
    WHERE 
        (c.sender = ? AND c.reciever = ?) 
        OR (c.sender = ? AND c.reciever = ?)
    `

	err := db.QueryRow(existingConversation,
		currentUser, currentUser, // Pour les CASE
		currentUser, user2, // Premier cas
		user2, currentUser, // Deuxième cas
	).Scan(&conversationID, &otherUserPicture, &otherUserUsername)

	if err == nil {
		return &Conversation{
			ConversationID:   conversationID,
			User1ID:          currentUser,
			User2ID:          user2,
			ReceiverPicture:  otherUserPicture,
			ReceiverUsername: otherUserUsername,
		}, nil
	} else if err != sql.ErrNoRows {
		return nil, fmt.Errorf("erreur lors de la vérification de la conversation: %v", err)
	}

	// Si aucune conversation n'existe, en créer une nouvelle
	conversationUUID, err := posts.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la génération du UUID: %v", err)
	}

	// Création de la nouvelle conversation
	creationDate := time.Now()
	createConversationQuery := `INSERT INTO conversations (conversation_uuid, sender, reciever, created_at) VALUES (?,?,?,?)`
	_, err = db.Exec(createConversationQuery, conversationUUID, currentUser, user2, creationDate)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création de la conversation: %v", err)
	}

	// Récupérer les infos de l'autre utilisateur
	err = db.QueryRow(`SELECT profile_picture, username FROM users WHERE user_uuid = ?`, user2).Scan(&otherUserPicture, &otherUserUsername)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des informations de l'autre utilisateur: %v", err)
	}

	return &Conversation{
		ConversationID:   conversationUUID,
		User1ID:          currentUser,
		User2ID:          user2,
		CreatedAt:        creationDate,
		ReceiverPicture:  otherUserPicture,
		ReceiverUsername: otherUserUsername,
	}, nil
}

func GetConversations(db *sql.DB, user_uuid string) ([]Conversation, error) {
	getConversations := `
    SELECT 
        c.conversation_uuid, 
        c.created_at,
        CASE 
            WHEN c.sender = ? THEN c.reciever
            ELSE c.sender
        END AS other_user_uuid,
        CASE 
            WHEN c.sender = ? THEN u2.profile_picture
            ELSE u1.profile_picture
        END AS other_user_profile_picture,
        CASE 
            WHEN c.sender = ? THEN u2.username
            ELSE u1.username
        END AS other_user_username
    FROM 
        conversations c
    JOIN users u1 ON u1.user_uuid = c.sender
    JOIN users u2 ON u2.user_uuid = c.reciever
    WHERE 
        c.sender = ? OR c.reciever = ?;
    `

	rows, err := db.Query(getConversations, user_uuid, user_uuid, user_uuid, user_uuid, user_uuid)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des conversations: %v", err)
	}
	defer rows.Close()

	var conversations []Conversation
	for rows.Next() {
		var conversation Conversation
		err := rows.Scan(
			&conversation.ConversationID,
			&conversation.CreatedAt,
			&conversation.User2ID,
			&conversation.ReceiverPicture,
			&conversation.ReceiverUsername,
		)
		if err != nil {
			return nil, fmt.Errorf("erreur lors du scan des conversations: %v", err)
		}
		conversations = append(conversations, conversation)
	}

	return conversations, nil
}
