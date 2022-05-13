package chat

import (
	"fmt"
	"log"
	"net/http"
	"stocks-chat/chat/utils"
	"stocks-chat/tools"
	"strings"

	"github.com/gorilla/websocket"
)

type Chat struct {
	users    map[string]*UserChat
	messages chan *Message
	userIn   chan *UserChat
	userOut  chan *UserChat
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return r.Method == http.MethodGet
	},
}

func (c *Chat) Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("Websocket connection error:", err.Error())
	}

	keys := r.URL.Query()
	username := keys.Get("username")
	if strings.TrimSpace(username) == "" {
		username = fmt.Sprintf("anon-%d", utils.GetRandomI64())
	}

	user := &UserChat{
		Username: username,
		Conn:     conn,
		Global:   c,
	}
	c.userIn <- user
	user.Read()
}

func (c *Chat) Run() {
	for {
		select {
		case user := <-c.userIn:
			c.addUserChat(user)
		case message := <-c.messages:
			c.broadcast(message)
		case user := <-c.userOut:
			c.disconnect(user)
		}
	}
}

func (c *Chat) addUserChat(userChat *UserChat) {
	if _, ok := c.users[userChat.Username]; !ok {
		c.users[userChat.Username] = userChat
		body := fmt.Sprintf("%s in to the room", userChat.Username)
		c.broadcast(NewMessage(body, "Server"))
	}
}

func (c *Chat) broadcast(message *Message) {
	for _, user := range c.users {
		user.Write(message)
	}

	if strings.Contains(message.Body, tools.GetDotEnvVariable("STOCK_COMMAND", "/stock=")) {
		messageSplit := strings.Split(message.Body, "=")
		callBot(messageSplit[1], c)
	}
}

func (c *Chat) disconnect(user *UserChat) {
	if _, ok := c.users[user.Username]; ok {
		defer user.Conn.Close()
		delete(c.users, user.Username)

		body := fmt.Sprintf("%s left the chat", user.Username)
		c.broadcast(NewMessage(body, "Server"))
	}
}

func Start(port string) {

	log.Printf("Chat running on http://localhost%s\n", port)
	c := &Chat{
		users:    make(map[string]*UserChat),
		messages: make(chan *Message),
		userIn:   make(chan *UserChat),
		userOut:  make(chan *UserChat),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Stocks chat!"))
	})
	http.HandleFunc("/chat", c.Handler)

	go c.Run()
	log.Fatal(http.ListenAndServe(port, nil))
}
