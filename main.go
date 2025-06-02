package main

import (
	"log"
	"os"
	"time"

	"github.com/Flop4ik/telegram-asr/packages/handlers"

	"github.com/joho/godotenv"

	tg "gopkg.in/telebot.v4"
)

func main() {

	godotenv.Load()

	pref := tg.Settings{
		Token:  os.Getenv("TELEGRAM_TOKEN"),
		Poller: &tg.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tg.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle(tg.OnVoice, func(c tg.Context) error {
		return handlers.OnVoice(c, b)
	})

	b.Handle("/start", func(c tg.Context) error {
		return handlers.StartCommand(c)
	})

	b.Start()

}
