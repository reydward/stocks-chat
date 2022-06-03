package stocks

import (
	"fmt"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"stocks-chat/tools"
	"strings"
)

func GetStockfromAPI(ticker string) bool {
	apiUrl := tools.GetDotEnvVariable("STOCKS_URL", "STOCKS_URL=https://stooq.com/q/l/?f=sd2t2ohlcv&h&e=csv&s=") + ticker
	response, err := http.Get(apiUrl)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return SendToQueue(buildResponse(string(responseData)))
}

func SendToQueue(message string) bool {
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

	q, err := ch.QueueDeclare("chat-messages", false, false, false, false, nil)

	if err != nil {
		log.Fatal(err)
	}

	err = ch.Publish("", q.Name, false, false,
		amqp.Publishing{
			Headers:     nil,
			ContentType: "text/plain",
			Body:        []byte(message),
		})

	if err != nil {
		return false
	}

	fmt.Println("publish: " + message)
	return true

}

func buildResponse(responseData string) string {
	splitResponseData := strings.Split(responseData, ",")
	ticker := splitResponseData[7]
	tickerFmt := ticker[8:]
	return tickerFmt + " quote is " + splitResponseData[13] + " per share"
}
