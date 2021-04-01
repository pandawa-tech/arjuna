package main

import (
    "fmt"
    "log"
    "net/http"
	"os"

    "github.com/line/line-bot-sdk-go/linebot"
    "github.com/joho/godotenv"
    "github.com/Maldris/mathparse"
)

var bot *linebot.Client

func defaultHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love you!")
}

func lineHandler(w http.ResponseWriter, r *http.Request) {
    events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
                expression := message.Text
                p := mathparse.NewParser(expression)
                var result string
                p.Resolve()
                if p.FoundResult() {
                    var resultFloat float64
                    resultFloat = p.GetValueResult()
                    result = fmt.Sprintf("%g", resultFloat)
                } else {
                    result = fmt.Sprintf("Gw gtau mksud cht anda: %s", p.GetExpressionResult())
                }
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(result)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}

func main() {
    var err error
    error := godotenv.Load(".env")
    if error != nil {
        log.Println(error)
	}
    bot, err = linebot.New(os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_ACCESS_TOKEN"))
    if err != nil {
        log.Println(err)
		log.Fatal("Failed initialize line chatbot")
	}
	fmt.Fprintf(os.Stdout, "Web Server started. Listening on port 8080\n")
    http.HandleFunc("/webhooks/line", lineHandler)
    http.HandleFunc("/", defaultHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}