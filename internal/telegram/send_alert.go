package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func SendTeleAlert(bot *tgbotapi.BotAPI, teleChatID int64, textCh <-chan string, errCh chan<- error) {
	for {
		msg := tgbotapi.NewMessage(teleChatID, <-textCh)
		msg.DisableWebPagePreview = true
		msg.ParseMode = "markdown"
		_, err := bot.Send(msg)
		if err != nil {
			errCh <- err
		}
	}
}
