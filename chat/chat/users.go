package chat

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type UserChat struct {
	Username string
	Conn     *websocket.Conn
	Global   *Chat
}

func (u *UserChat) Read() {
	for {
		if _, message, err := u.Conn.ReadMessage(); err != nil {
			log.Println("Error on read message:", err.Error())

			break
		} else {
			u.Global.messages <- NewMessage(string(message), u.Username)
		}
	}

	u.Global.userOut <- u
}

func (u *UserChat) Write(message *Message) {
	b, _ := json.Marshal(message)

	if err := u.Conn.WriteMessage(websocket.TextMessage, b); err != nil {
		log.Println("Error on write message:", err.Error())
	}
}
