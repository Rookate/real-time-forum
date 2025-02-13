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
	sender, _ := authentification.GetUserFromCookie(r)
	receiver := user2
	var receiverPicture, receiverUsername string

	// Vérifier si une conversation existe déjà entre les deux utilisateurs
	existingConversation := `
	SELECT 
	    c.conversation_uuid, 
	    u.profile_picture, 
	    u.username
	FROM 
	    conversations c
	JOIN 
	    users u ON u.user_uuid = c.reciever
	WHERE 
	    (c.sender = ? AND c.reciever = ?) 
	    OR 
	    (c.sender = ? AND c.reciever = ?);
	`

	var conversationID string
	err := db.QueryRow(existingConversation, sender, receiver, sender, receiver).Scan(&conversationID, &receiverPicture, &receiverUsername)
	if err == nil {
		// Une conversation existe déjà, la retourner sans en créer une nouvelle
		return &Conversation{
			ConversationID:   conversationID,
			User1ID:          sender,
			User2ID:          receiver,
			ReceiverUsername: receiverUsername,
			ReceiverPicture:  receiverPicture,
		}, nil
	} else if err != sql.ErrNoRows {
		return nil, fmt.Errorf("erreur lors de la vérification de la conversation: %v", err)
	}

	// Si aucune conversation n'existe, en créer une nouvelle
	conversationUUID, err := posts.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la génération du UUID: %v", err)
	}

	// Insertion dans la base de données
	creationDate := time.Now()
	createConversationQuery := `INSERT INTO conversations (conversation_uuid, sender, reciever, created_at) VALUES (?,?,?,?)`
	_, err = db.Exec(createConversationQuery, conversationUUID, sender, receiver, creationDate)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création de la conversation: %v", err)
	}

	// Récupérer la photo de profil et le nom d'utilisateur du receiver
	err = db.QueryRow(`SELECT profile_picture, username FROM users WHERE user_uuid = ?`, receiver).Scan(&receiverPicture, &receiverUsername)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des informations du receiver: %v", err)
	}

	// Créer et retourner la nouvelle conversation
	newConversation := &Conversation{
		ConversationID:   conversationUUID,
		User1ID:          sender,
		User2ID:          receiver,
		CreatedAt:        creationDate,
		ReceiverUsername: receiverUsername,
		ReceiverPicture:  receiverPicture,
	}

	return newConversation, nil
}
