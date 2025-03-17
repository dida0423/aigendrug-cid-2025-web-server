package chat

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/gorilla/websocket"
	"github.com/openai/openai-go"
)

/*
WebSocketHandler is a handler function for Gin routes that upgrades the HTTP connection to a WebSocket connection.
It reads messages from the client and broadcast them to all other clients in the same session.

HandleMessages is a function that reads messages from the broadcast channel and sends them to all clients in the session.
If the message is from the user, it generates an AI response and sends it to all clients.
HandleMessages is called as a goroutine in the main function and runs in the context of the main gin server.

WebSocketHandler and HandleMessages use a shared map called sessionClients to manage WebSocket clients in the same session.
*/

// Websocket Upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var sessionClients = make(map[string]map[*websocket.Conn]bool) // Dynamic map of sessionID to clients
var broadcast = make(chan ChatMessage)                         // Sync channel for broadcasting messages
var mutex = sync.Mutex{}                                       // Mutex for sessionClients

// Replace with AgentResponse from aigendrug ai service
func generateAIResponse(message string) string {
	client := openai.NewClient()
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(message),
		}),
		Model: openai.F(openai.ChatModelGPT4o),
	})
	if err != nil {
		panic(err.Error())
	}

	return chatCompletion.Choices[0].Message.Content
}

func WebSocketHandler(c *gin.Context, db *gocql.Session) {
	sessionID := c.Query("sessionID")
	if sessionID == "" {
		c.JSON(400, gin.H{"error": "sessionID is required"})
		return
	}

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket Upgrade Error:", err)
		return
	}
	defer conn.Close()

	// Lock mutex of sessionClients and add the new client.
	mutex.Lock()
	if sessionClients[sessionID] == nil {
		sessionClients[sessionID] = make(map[*websocket.Conn]bool)
	}
	sessionClients[sessionID][conn] = true
	mutex.Unlock()

	// Read messages from the client and broadcast them to all other clients
	for {
		var msg CreateChatMessageDTO
		err := conn.ReadJSON(&msg)
		// If there is an error, remove the client from the sessionClients map and break the loop
		if err != nil {
			log.Println("Read Error:", err)
			mutex.Lock()
			delete(sessionClients[sessionID], conn)
			mutex.Unlock()
			break
		}

		// Save the message to the database and broadcast it to all clients
		err = saveChatMessageToDB(db, &msg)
		if err != nil {
			log.Println("DB Save Error:", err)
			continue
		}

		broadcast <- ChatMessage{
			SessionID:     msg.SessionID,
			Role:          msg.Role,
			Message:       msg.Message,
			MessageType:   msg.MessageType,
			LinkedToolIDs: msg.LinkedToolIDs,
		}
	}
}

func saveChatMessageToDB(db *gocql.Session, msg *CreateChatMessageDTO) error {
	newUUID, err := gocql.RandomUUID()
	if err != nil {
		return err
	}

	query := db.Query("INSERT INTO chat_messages (id, session_id, role, message, created_at, message_type, linked_tool_ids) VALUES (?, ?, ?, ?, toTimestamp(now()), ?, ?)",
		newUUID, msg.SessionID, msg.Role, msg.Message, msg.MessageType, msg.LinkedToolIDs)

	if err := query.Exec(); err != nil {
		return err
	}

	return nil
}

// Read messages from the broadcast channel and send them to all clients in the session
func HandleMessages(db *gocql.Session) {
	for {
		msg := <-broadcast

		// Lock the mutex and check if the session exists
		mutex.Lock()
		clients, exists := sessionClients[msg.SessionID.String()]
		if !exists {
			mutex.Unlock()
			continue
		}

		// For each client in the session, send the message
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println("Send Error:", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()

		// If the message is from the user, generate an AI response and send it to all clients
		if msg.Role == ChatRoleUser {
			aiResponse := generateAIResponse(msg.Message)

			aiMsg := ChatMessage{
				SessionID:     msg.SessionID,
				Role:          ChatRoleAssistant,
				Message:       aiResponse,
				MessageType:   msg.MessageType,
				LinkedToolIDs: msg.LinkedToolIDs,
			}
			err := saveChatMessageToDB(db, &CreateChatMessageDTO{
				SessionID:     aiMsg.SessionID,
				Role:          aiMsg.Role,
				Message:       aiMsg.Message,
				MessageType:   aiMsg.MessageType,
				LinkedToolIDs: aiMsg.LinkedToolIDs,
			})
			if err != nil {
				log.Println("Failed to save AI response:", err)
				continue
			}

			// finishMsg is a message to indicate that the AI has finished responding.
			// Web clients can use this message to manage the chat UI.(e.g., stop loading animation, scroll to the bottom of the chat)
			finishMsg := map[string]interface{}{
				"status": "finished",
			}

			// Lock the mutex and send the AI response and finish message to all clients
			mutex.Lock()
			for client := range clients {
				err := client.WriteJSON(aiMsg)
				if err != nil {
					log.Println("Send Error:", err)
					client.Close()
					delete(clients, client)
				}

				err = client.WriteJSON(finishMsg)
				if err != nil {
					log.Println("Send Error:", err)
					client.Close()
					delete(clients, client)
				}
			}
			mutex.Unlock()
		}
	}
}
