package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Flop4ik/telegram-asr/packages/gemini"

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
		messageID := strconv.Itoa(int(c.Message().ID))

		currentDir, err := os.Getwd()
		if err != nil {
			log.Printf("Failed to get current working directory: %v", err)
			return c.Send("Failed to get current working directory.")
		}

		path := filepath.Join(currentDir, "tmp-voices", messageID+".ogg")

		log.Printf("Received voice message from %s", c.Sender().Username)

		b.Download(&c.Message().Voice.File, path)

		fmt.Println(path)

		result, err := gemini.RecognizeText(path)

		if err != nil {
			log.Printf("Error recognizing text: %v", err)
			return c.Send("Error recognizing text from the voice message.")
		}

		fmt.Println(path)
		if err := os.Remove(path); err != nil {
			log.Printf("Failed to delete file %s: %v", path, err)
		}

		return c.Send(result)
	})

	b.Start()

}
