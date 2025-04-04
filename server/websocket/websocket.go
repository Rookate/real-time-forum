package ws

import (
	"encoding/json"
	"fmt"
	"forum/server"
	authentification "forum/server/api/login"
	"forum/server/message"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Structure pour stocker les informations du client
type Client struct {
	conn             *websocket.Conn
	conversationUUID string
	user_UUID        string
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Modifier la map pour stocker les informations du client
var clients = make(map[*websocket.Conn]*Client)
var broadcast = make(chan []byte)
var mutex = &sync.Mutex{}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}
	defer conn.Close()

	// Créer un nouveau client
	client := &Client{
		conn:             conn,
		conversationUUID: "",
		user_UUID:        "",
	}

	mutex.Lock()
	clients[conn] = client
	user_uuid, err := authentification.GetUserFromCookie(r)
	if err == nil {
		client.user_UUID = user_uuid
	}
	clients[conn].user_UUID = user_uuid
	mutex.Unlock()
	broadcastActiveUsers()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			break
		}

		var request map[string]interface{}
		if err := json.Unmarshal(msg, &request); err != nil {
			fmt.Println("Error parsing message:", err)
			continue
		}

		messageType, _ := request["type"].(string)

		switch messageType {
		case "single_message":
			content, ok := request["content"].(map[string]interface{})
			if !ok {
				fmt.Println("Invalid message content")
				continue
			}

			// Récupérer l'UUID de la conversation du message
			conversationUUID := content["conversation_uuid"].(string)

			response, _ := json.Marshal(map[string]interface{}{
				"type":                   "single_message",
				"sender_username":        content["sender_username"],
				"sender_profile_picture": content["sender_profile_picture"],
				"content":                content["content"],
				"created_at":             time.Now().Format(time.RFC1123),
			})

			// Envoyer uniquement aux clients dans la même conversation
			mutex.Lock()
			for _, client := range clients {
				if client.conversationUUID == conversationUUID {
					err := client.conn.WriteMessage(websocket.TextMessage, response)
					if err != nil {
						client.conn.Close()
						delete(clients, client.conn)
					}
				}
			}
			mutex.Unlock()

		case "getMessages":
			conversationUUID, ok := request["conversation_uuid"].(string)
			if !ok || conversationUUID == "" {
				fmt.Println("Invalid conversation_uuid")
				continue
			}

			offset, _ := request["offset"].(float64)
			// if (offsetOk) {

			// }
			limit, _ := request["limit"].(float64)
			//if (limitOk) {}

			// Mettre à jour l'UUID de la conversation pour ce client
			mutex.Lock()
			clients[conn].conversationUUID = conversationUUID
			mutex.Unlock()

			// log.Printf("Data pour recevoir les messages : conversation : %s, offset :%f, limit : %f", conversationUUID, offset, limit)

			messages, err := message.GetMessagesByConversations(server.Db, r, conversationUUID, int(offset), int(limit))
			if err != nil {
				fmt.Println("Error fetching messages:", err)
				continue
			}

			response, _ := json.Marshal(map[string]interface{}{
				"type":     "messages",
				"messages": messages,
			})

			// Envoyer uniquement au client qui a demandé les messages
			conn.WriteMessage(websocket.TextMessage, response)
		case "getMoreMessages":
			conversationUUID, ok := request["conversation_uuid"].(string)
			if !ok || conversationUUID == "" {
				fmt.Println("Invalid conversation_uuid")
				continue
			}

			offset, _ := request["offset"].(float64)
			limit, _ := request["limit"].(float64)
			mutex.Lock()
			clients[conn].conversationUUID = conversationUUID
			mutex.Unlock()

			messages, err := message.GetMessagesByConversations(server.Db, r, conversationUUID, int(offset), int(limit))
			if err != nil {
				fmt.Println("Error fetching messages:", err)
				continue
			}

			// log.Printf("Data pour recevoir PLUS de messages : conversation : %s, offset :%f, limit : %f", conversationUUID, offset, limit)

			response, _ := json.Marshal(map[string]interface{}{
				"type":     "moreMessages",
				"messages": messages,
			})

			// Envoyer uniquement au client qui a demandé les messages
			conn.WriteMessage(websocket.TextMessage, response)

		case "typing":
			typing, _ := request["isTyping"].(bool)
			response, _ := json.Marshal(map[string]interface{}{
				"type":     "typing",
				"isTyping": typing,
			})

			sender := client

			for _, client := range clients {
				if client != sender {
					err := client.conn.WriteMessage(websocket.TextMessage, response)
					if err != nil {
						client.conn.Close()
						delete(clients, client.conn)
					}
				}
			}
		case "notification":
			fmt.Println("Notification")
			response, _ := json.Marshal(map[string]interface{}{
				"type": "notification",
			})

			// On envoi a l'autre utilisateur
			sender := client

			for _, client := range clients {
				if client != sender {
					err := client.conn.WriteMessage(websocket.TextMessage, response)
					if err != nil {
						client.conn.Close()
						delete(clients, client.conn)
					}
				}

			}
		}
	}
}

func HandleMessages() {
	for {
		message := <-broadcast
		mutex.Lock()
		for _, client := range clients {
			err := client.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				client.conn.Close()
				delete(clients, client.conn)
			}
		}
		mutex.Unlock()
	}
}

func broadcastActiveUsers() {
	mutex.Lock()
	defer mutex.Unlock()

	var activeUsers []string
	for _, client := range clients {
		if client.user_UUID != "" {
			activeUsers = append(activeUsers, client.user_UUID)
		}
	}

	response, _ := json.Marshal(map[string]interface{}{
		"type":         "active_users",
		"active_users": activeUsers,
	})

	// Envoyer la liste des utilisateurs à tous les clients connectés
	for _, client := range clients {
		err := client.conn.WriteMessage(websocket.TextMessage, response)
		if err != nil {
			client.conn.Close()
			delete(clients, client.conn)
		}
	}
}
