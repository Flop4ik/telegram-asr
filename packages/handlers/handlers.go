package handlers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Flop4ik/telegram-asr/packages/gemini"

	tg "gopkg.in/telebot.v4"
)

func OnVoice(c tg.Context, b *tg.Bot) error {

	if c.Message().Voice.Duration > 1200 {
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

	result, err := gemini.RecognizeText(path)

	if err != nil {
		log.Printf("Error recognizing text: %v", err)
		return c.Send("Error recognizing text from the voice message.")
	}

	fmt.Println(path)
	if err := os.Remove(path); err != nil {
		log.Printf("Failed to delete file %s: %v", path, err)
	}

	return c.Send(result, &tg.SendOptions{ParseMode: tg.ModeMarkdown})
}

func StartCommand(c tg.Context) error {
	log.Printf("Received /start command from %s", c.Sender().Username)
	welcomeMessage := `
👋 *Добро пожаловать в TASR !* 🎤

Что такое TASR? Это аббревиатура от "Telegram ASR" (AI Speech Recognition). 🤖

Я здесь, чтобы помочь вам преобразовать ваши голосовые сообщения в текст. 📝

Этот бот нужен, чтобы переводить голосовые в текст, без telegram premium

Просто отправьте мне голосовое сообщение, и я сделаю все возможное, чтобы ее расшифровать. ✨
`
	return c.Send(welcomeMessage, &tg.SendOptions{ParseMode: tg.ModeMarkdown})
}
