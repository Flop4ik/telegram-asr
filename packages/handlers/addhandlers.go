package handlers

import (
	"github.com/Flop4ik/telegram-asr/packages/handlers/actions"
	coms "github.com/Flop4ik/telegram-asr/packages/handlers/commands"
	tg "gopkg.in/telebot.v4"
)

func AddHandlers(b *tg.Bot) {
	b.Handle("/start", func(c tg.Context) error {
		return coms.StartCommand(c)
	})
	b.Handle("/help", func(c tg.Context) error {
		return coms.HelpCommand(c)

	})
	b.Handle("/tokens", func(c tg.Context) error {
		return coms.CheckTokens(c)
	})
	b.Handle("/changemode", func(c tg.Context) error {
		return coms.ChangeMode(c)
	})
	b.Handle("/mode", func(c tg.Context) error {
		return coms.CheckMode(c)
	})

	b.Handle(tg.OnVoice, func(c tg.Context) error {
		return actions.OnVoice(c, b)
	})

}
