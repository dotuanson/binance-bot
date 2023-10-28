package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func SendTeleAlert(bot *tgbotapi.BotAPI, teleChatID int64, textCh <-chan string, errCh chan<- error) {
	for {
		text := tgbotapi.NewMessage(teleChatID, <-textCh)
		_, err := bot.Send(text)
		if err != nil {
			errCh <- err
		}
	}
}
