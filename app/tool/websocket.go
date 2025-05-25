package tool

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openai/openai-go"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var sessionClients = make(map[string]map[*websocket.Conn]bool) // Dynamic map of sessionID to clients
var broadcast = make(chan ToolMessage)                         // Sync channel for broadcasting messages
var mutex = sync.Mutex{}                                       // Mutex for sessionClients

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

func WebSocketHandler(c *gin.Context, db *pgxpool.Pool) {
	sessionID := c.Query("sessionID")
	if sessionID == "" {
		c.JSON(400, gin.H{"error": "sessionID is required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket Upgrade Error:", err)
		return
	}
	defer conn.Close()

	mutex.Lock()
	if sessionClients[sessionID] == nil {
		sessionClients[sessionID] = make(map[*websocket.Conn]bool)
	}
	sessionClients[sessionID][conn] = true
	mutex.Unlock()

	for {
		var msg CreateToolMessageDTO
		err := conn.ReadJSON(&msg)

		if err != nil {
			log.Println("Read Error:", err)
			mutex.Lock()
			delete(sessionClients[sessionID], conn)
			mutex.Unlock()
			break
		}

		err = saveChatMessageToDB(db, &msg)
		if err != nil {
			log.Println("DB Save Error:", err)
			continue
		}

		broadcast <- ToolMessage{
			SessionID: msg.SessionID,
			ToolID:    msg.ToolID,
			Role:      msg.Role,
			Data:      msg.Data,
		}
	}
}

func saveChatMessageToDB(db *pgxpool.Pool, msg *CreateToolMessageDTO) error {
	newUUID := uuid.New()

	dataStr, err := json.Marshal(msg.Data)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		context.Background(),
		`INSERT INTO tool_messages (id, session_id, tool_id, role, data, created_at)
         VALUES ($1, $2, $3, $4, $5, $6)`,
		newUUID,
		msg.SessionID,
		msg.ToolID,
		msg.Role,
		string(dataStr),
		time.Now(),
	)
	return err
}

func HandleMessages(db *pgxpool.Pool) {
	for {
		msg := <-broadcast

		mutex.Lock()
		clients, exists := sessionClients[msg.SessionID.String()]
		if !exists {
			mutex.Unlock()
			continue
		}

		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println("Send Error:", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()

		if msg.Role == ToolRoleUser {
			aiResponse := generateAIResponse(msg.Data["message"].(string))

			aiMsg := ToolMessage{
				SessionID: msg.SessionID,
				ToolID:    msg.ToolID,
				Role:      ToolRoleAssistant,
				Data:      map[string]interface{}{"message": aiResponse},
				CreatedAt: msg.CreatedAt,
			}
			err := saveChatMessageToDB(db, &CreateToolMessageDTO{
				SessionID: msg.SessionID,
				ToolID:    msg.ToolID,
				Role:      ToolRoleAssistant,
				Data:      map[string]interface{}{"message": aiResponse},
			})
			if err != nil {
				log.Println("Failed to save AI response:", err)
				continue
			}

			finishMsg := map[string]interface{}{
				"status": "finished",
			}

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
