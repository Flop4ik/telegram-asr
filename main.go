package main

import (
	"log"
	"os"
	"time"

	db "github.com/Flop4ik/telegram-asr/packages/database"
	"github.com/Flop4ik/telegram-asr/packages/handlers"

	"github.com/joho/godotenv"

	tg "gopkg.in/telebot.v4"
)

func main() {

	currentDay := time.Now().Format("2006-01-02")

	godotenv.Load()
	if err := db.Initialize("./database.db"); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}

	pref := tg.Settings{
		Token:  os.Getenv("TELEGRAM_TOKEN"),
		Poller: &tg.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tg.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	handlers.AddHandlers(b)

	go func() {
		for {
			if time.Now().Format("2006-01-02") != currentDay {
				currentDay = time.Now().Format("2006-01-02")
				log.Println("New day detected, resetting tokens for all users.")
				err := db.ResetTokens()
				if err != nil {
					log.Printf("Failed to reset tokens: %v", err)
				} else {
					log.Println("Tokens reset successfully.")
				}
			}
			time.Sleep(10 * time.Minute)
		}
	}()

	b.Start()

}
