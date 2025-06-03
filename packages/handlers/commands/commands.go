package handlers

import (
	"fmt"
	"log"

	db "github.com/Flop4ik/telegram-asr/packages/database"
	tg "gopkg.in/telebot.v4"
)

func StartCommand(c tg.Context) error {
	log.Printf("Received /start command from %s", c.Sender().Username)
	welcomeMessage := `
👋 *Добро пожаловать в TASR !* 🎤

Что такое TASR? Это аббревиатура от "Telegram AI Speech Recognition. 🤖

Я здесь, чтобы помочь вам преобразовать ваши голосовые сообщения в текст. 📝

Этот бот нужен, чтобы переводить голосовые в текст, без telegram premium.

Просто отправьте мне голосовое сообщение, и я сделаю все возможное, чтобы ее расшифровать. ✨

*Чтобы узнать, как пользоваться ботом, используйте команду /help.*
`

	db.CreateUser(c.Sender().ID)
	return c.Send(welcomeMessage, &tg.SendOptions{ParseMode: tg.ModeMarkdown})
}

func HelpCommand(c tg.Context) error {
	log.Printf("Received /help command from %s", c.Sender().Username)
	helpMessage := `
👋 *Добро пожаловать в TASR!*

Я помогу тебе превратить голосовые сообщения в текст! 📝

📌 *Как пользоваться:*
1. Просто перешли мне голосовое сообщение.
2. Я его расшифрую и пришлю тебе текст.

⚙️ *Режимы работы:*
- По умолчанию я просто перевожу голос в текст (транскрибация).
- Ты можешь изменить режим командой /changemode.
- Узнать текущий режим можно командой /mode.

💬 *Как посмотреть токены?*
- Чтобы узнать, сколько токенов у тебя осталось, используй команду /tokens.

*Доступные режимы:*
1.  🎯 Транскрибация (по умолчанию): Полный перевод аудио в текст.
    Стоимость: 1 токен за запрос.
2.  📋 Краткий пересказ: Я сделаю краткое содержание аудио.
    Стоимость: 1.2 токена за запрос.

🪙 *Токены:*
- У тебя есть *150 токенов* в день для использования моих функций.
- Режим краткого пересказа тратит 15 токенов за 1 аудио.
- Режим транскрибации тратит 10 токенов за 1 аудио.
- Токены обновляются раз в день.
- Токены обновляются ежедневно.

✨ *Просто отправь голосовое, и я начну работу!*
`
	return c.Send(helpMessage, &tg.SendOptions{ParseMode: tg.ModeMarkdown})
}

func CheckTokens(c tg.Context) error {
	id := c.Sender().ID
	tokens, err := db.GetTokens(id)
	if err != nil {
		log.Printf("Failed to get tokens for user %d: %v", id, err)
		return c.Send("❌ Ошибка при получении токенов. Пожалуйста, попробуйте позже.")
	}

	message := `
🪙 *Доступные токены*

У вас осталось *%d из 150* токенов.

Токены обновляются автоматически каждый день.
`
	return c.Send(fmt.Sprintf(message, tokens), &tg.SendOptions{ParseMode: tg.ModeMarkdown})
}

func CheckMode(c tg.Context) error {
	id := c.Sender().ID
	mode, err := db.GetMode(id)
	if err != nil {
		log.Printf("Failed to get mode for user %d: %v", id, err)
		return c.Send("❌ Ошибка при получении режима работы. Пожалуйста, попробуйте позже.")
	}
	switch mode {
	case "transcribe":
		mode = "Транскрибация"
	case "summarize":
		mode = "Краткий пересказ"
	default:
		log.Printf("Unknown mode for user %d: %s", id, mode)
		return c.Send("❌ Неизвестный режим работы. Пожалуйста, попробуйте позже.")
	}
	message := `
🛠️ *Текущий режим работы*
Ваш текущий режим работы: *%s*.
Чтобы изменить режим, используйте команду /changemode.
`
	return c.Send(fmt.Sprintf(message, mode), &tg.SendOptions{ParseMode: tg.ModeMarkdown})
}

func ChangeMode(c tg.Context) error {
	id := c.Sender().ID
	mode, err := db.GetMode(id)
	if err != nil {
		log.Printf("Failed to get mode for user %d: %v", id, err)
		return c.Send("❌ Ошибка при получении режима работы. Пожалуйста, попробуйте позже.")
	}
	switch mode {
	case "transcribe":
		err = db.SetMode(id, "summarize")
		if err != nil {
			log.Printf("Failed to set mode for user %d: %v", id, err)
			return c.Send("❌ Ошибка при смене режима. Пожалуйста, попробуйте позже.")
		}
		return c.Send("✅ Режим работы изменен на *Краткий пересказ*.", &tg.SendOptions{ParseMode: tg.ModeMarkdown})
	case "summarize":
		err = db.SetMode(id, "transcribe")
		if err != nil {
			log.Printf("Failed to set mode for user %d: %v", id, err)
			return c.Send("❌ Ошибка при смене режима. Пожалуйста, попробуйте позже.")
		}
		return c.Send("✅ Режим работы изменен на *Транскрибация*.", &tg.SendOptions{ParseMode: tg.ModeMarkdown})
	default:
		log.Printf("Unknown mode for user %d: %s", id, mode)
		return c.Send("❌ Неизвестный режим работы. Пожалуйста, попробуйте позже.")
	}
}
