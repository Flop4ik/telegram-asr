package actions

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	db "github.com/Flop4ik/telegram-asr/packages/database"
	"github.com/Flop4ik/telegram-asr/packages/gemini"

	tg "gopkg.in/telebot.v4"
)

func OnVoice(c tg.Context, b *tg.Bot) error {

	var requiredTokens int32

	id := c.Sender().ID

	mode, err := db.GetMode(id)

	if err != nil {
		log.Printf("Failed to get mode for user %d: %v", id, err)
		return c.Send("Ошибка при получении режима работы. Пожалуйста, попробуйте позже.")
	}

	tokens, err := db.GetTokens(id)
	if err != nil {
		log.Printf("Failed to get tokens for user %d: %v", id, err)
		return c.Send("Ошибка при получении токенов. Пожалуйста, попробуйте позже.")
	}

	switch mode {
	case "transcribe":
		requiredTokens = 10
	case "summarize":
		requiredTokens = 15
	default:
		log.Printf("Unknown mode for user %d: %s", id, mode)
		requiredTokens = 10
	}

	if tokens < requiredTokens {
		return c.Send("У вас недостаточно токенов для транскрибации. Пожалуйста, попробуйте завтра.")
	}

	if c.Message().Voice.Duration > 600 {
		log.Printf("Voice message from %s exceeds 10 minutes, ignoring.", c.Sender().Username)
		return c.Send("Голосовое сообщение слишком длинное. Пожалуйста, отправьте голосовое сообщение длиной до 10 минут.")
	}

	fmt.Println("Received voice message")
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

	result, err := gemini.RecognizeText(path, mode)

	if err != nil {
		log.Printf("Error recognizing text: %v", err)
		return c.Send("Error recognizing text from the voice message.")
	}

	fmt.Println(path)
	if err := os.Remove(path); err != nil {
		log.Printf("Failed to delete file %s: %v", path, err)
	}

	db.RemoveTokens(c.Sender().ID)

	tokens, _ = db.GetTokens(id)

	c.Send(result, &tg.SendOptions{ParseMode: tg.ModeMarkdown})
	return c.Send(fmt.Sprintf("🪙 У вас осталось *%d из 150* токенов.", tokens), &tg.SendOptions{ParseMode: tg.ModeMarkdown})
}
