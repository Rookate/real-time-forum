package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"websocket/server"
	apiConversations "websocket/server/api/conversations"
	apiMessages "websocket/server/api/message"
	"websocket/server/message"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool) // <--- Clients connectés
var broadcast = make(chan []byte)            // <--- Channels différents
var mutex = &sync.Mutex{}                    // <--- // Protect clients map

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}
	defer conn.Close()

	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

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
		if request["type"] == "getMessages" {
			conversationUUID, ok := request["conversation_uuid"].(string)
			if !ok || conversationUUID == "" {
				fmt.Println("Invalid conversation_uuid")
				continue
			}

			messages, err := message.GetMessagesByConversations(server.DB, r, conversationUUID)
			if err != nil {
				fmt.Println("Error fetching messages:", err)
				continue
			}

			response, _ := json.Marshal(map[string]interface{}{
				"type":     "messages",
				"messages": messages,
			})

			conn.WriteMessage(websocket.TextMessage, response)
		} else {
			// Rediriger les autres messages au channel broadcast
			broadcast <- msg
		}
	}
}
func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		message := <-broadcast

		// Send the message to all connected clients
		mutex.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
func conversationHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) > 2 {
		conversationUUID := parts[2]
		fmt.Printf("UUID de la conversation: %s\n", conversationUUID)
	}
	http.ServeFile(w, r, "conversation.html")
}

func main() {
	server.Test()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/ws", wsHandler)

	// API endpoints
	//http.HandleFunc("/api/post/createPost", post.CreatePostHandler)

	go handleMessages()
	http.HandleFunc("/", homeHandler)
	http.HandleFunc(`/conversation/`, conversationHandler)
	http.HandleFunc("/api/message/createMessage", apiMessages.CreateMessage)
	http.HandleFunc("/api/message/getMessage", apiMessages.GetMessageByConversation)
	fmt.Println("Welcome")
	http.HandleFunc("/api/conversations/createConversation", apiConversations.CreateConversation)
	fmt.Println("Welcome1")
	fmt.Println("WebSocket server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
