package server

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("sqlite3", "./data.db")
	fmt.Println("DATABASE")
	if err == nil {
		err = DB.Ping()
	}
	log.Printf("Tentative de connexion à la base de données échouée (%d/%d). Nouvelle tentative dans 5 secondes...")
	time.Sleep(5 * time.Second)

	if err != nil {
		log.Fatalf("Impossible de se connecter à la base de données après %d tentatives : %v", err)
	}
	createTables(DB)

	log.Println("Connexion à la base de données réussie")
}

func createTables(db *sql.DB) {
	// Requête pour créer la table users
	createUsersTable := `
    CREATE TABLE IF NOT EXISTS users (
        user_uuid TEXT PRIMARY KEY,
        username TEXT,
        profil_picture TEXT
    );`

	// Requête pour créer la table posts
	createConversationTable := `
  	CREATE TABLE IF NOT EXISTS conversations (
    conversation_uuid TEXT PRIMARY KEY,
    user1_uuid TEXT,
    user2_uuid TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user1_uuid) REFERENCES users(user_uuid),
    FOREIGN KEY (user2_uuid) REFERENCES users(user_uuid)
);
`

	createMessagesTable := `
	CREATE TABLE IF NOT EXISTS messages (
    message_uuid TEXT PRIMARY KEY,
    conversation_uuid TEXT,
    sender_uuid TEXT,
    receiver_uuid TEXT,
    message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (conversation_uuid) REFERENCES conversations(conversation_uuid),
    FOREIGN KEY (sender_uuid) REFERENCES users(user_uuid),
    FOREIGN KEY (receiver_uuid) REFERENCES users(user_uuid)
);
`
	// Exécution des requêtes pour créer les tables
	statements := []struct {
		name      string
		statement string
	}{
		{"users", createUsersTable},
		{"conversations", createConversationTable},
		{"messages", createMessagesTable},
	}

	var createdTables []string

	for _, stmt := range statements {
		_, err := db.Exec(stmt.statement)
		if err != nil {
			log.Fatalf("Erreur lors de la création de la table %s: %v", stmt.name, err)
		}
		// Ajoute le nom de la table créée
		createdTables = append(createdTables, stmt.name)
	}

	if len(createdTables) > 0 {
		fmt.Printf("Tables créées avec succès : %s\n", createdTables)
	} else {
		fmt.Println("Aucune table n'a été créée.")
	}
}
