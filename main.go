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

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love you!")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Healthy")
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
    port := getEnv("PORT", "8080")
    bot, err = linebot.New(os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_ACCESS_TOKEN"))
    if err != nil {
        log.Println(err)
		log.Fatal("Failed initialize line chatbot")
	}
    http.HandleFunc("/webhooks/line", lineHandler)
    http.HandleFunc("/healthcheck", healthCheckHandler)
    http.HandleFunc("/", defaultHandler)
    log.Print("Listening on port " + port + " ... ")
    err = http.ListenAndServe(":" + port, nil)
    if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}