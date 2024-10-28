package request

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/server"
	authentification "forum/server/api/login"
	posts "forum/server/utils"
	"net/http"
	"time"
)

type AdminRequest struct {
	Request_uuid    string    `json:"request_uuid"`
	User_uuid       string    `json:"user_uuid"`
	Created_at      time.Time `json:"created_at"`
	IsRead          bool      `json:"isRead"`
	Content         string    `json:"content"`
	Username        string    `json:"username"`
	Profile_Picture string    `json:"profile_picture"`
}

func CreateAdminRequest(db *sql.DB, r *http.Request, params map[string]interface{}) error {

	Request_uuid, _ := posts.GenerateUUID()

	user_UUID, err := authentification.GetUserFromCookie(r)
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération de l'UUID de l'utilisateur: %v", err)
	}

	content, contentOK := params["content"].(string)
	if !contentOK {
		return errors.New("contenu manquant")
	}

	createPostQuery := `INSERT INTO admin_requests (request_uuid, user_uuid, isRead, created_at, content) VALUES (?, ?, ?, ?, ?)`
	creationDate := time.Now()

	_, err = server.RunQuery(createPostQuery, Request_uuid, user_UUID, false, creationDate, content)
	if err != nil {
		return fmt.Errorf("erreur lors de la création de la demande: %v", err)
	}

	return nil
}

func FetchAdminRequest(db *sql.DB) ([]AdminRequest, error) {
	fetchRequestQuery := `
		SELECT ar.*, u.username, u.profile_picture
		FROM admin_requests ar
		JOIN users u ON ar.user_uuid = u.user_uuid
		WHERE ar.isRead = FALSE;`

	// Exécution de la requête avec server.RunQuery
	rows, err := server.RunQuery(fetchRequestQuery)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération de la demande: %v", err)
	}

	var adminRequests []AdminRequest

	// Boucle pour extraire les données des lignes
	for _, row := range rows {
		adminRequest := AdminRequest{}

		if val, ok := row["request_uuid"].(string); ok {
			adminRequest.Request_uuid = val
		}
		if val, ok := row["user_uuid"].(string); ok {
			adminRequest.User_uuid = val
		}
		if val, ok := row["created_at"].(time.Time); ok {
			adminRequest.Created_at = val
		}
		if val, ok := row["isRead"].(bool); ok {
			adminRequest.IsRead = val
		}
		if val, ok := row["content"].(string); ok {
			adminRequest.Content = val
		}

		if val, ok := row["username"].(string); ok {
			adminRequest.Username = val
		}

		if val, ok := row["profile_picture"].(string); ok {
			adminRequest.Profile_Picture = val
		}

		adminRequests = append(adminRequests, adminRequest)
	}

	return adminRequests, nil
}
