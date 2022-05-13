package chat

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"os"
	"stocks-chat/chat/utils"
	"stocks-chat/tools"
)

func callBot(ticker string, c *Chat) {
	botUrl := tools.GetDotEnvVariable("BOT_URL", "") + ticker
	response, err := http.Get(botUrl)

	if err != nil || response.Status != "200 OK" {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	go getMessagesFromQueue(c)
}

func getMessagesFromQueue(c *Chat) {
	conn, err := amqp.Dial(tools.GetDotEnvVariable("QUEUE_URL", "amqp://guest:guest@localhost:5672/"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	ch, err := conn.Channel()

	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()
	channelMessages, err := ch.Consume(
		"chat-messages",
		"",
		true,
		false,
		false,
		false, nil)

	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan bool)
	go func() {
		for channelMessage := range channelMessages {
			fmt.Println("consume: " + string(channelMessage.Body))
			stockMessage := &Message{
				ID:     utils.GetRandomI64(),
				Body:   string(channelMessage.Body),
				Sender: "Bot",
			}

			for _, user := range c.users {
				user.Write(stockMessage)
			}
		}
	}()
	<-forever
}
