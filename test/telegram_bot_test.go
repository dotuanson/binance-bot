package test

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSendAlert(t *testing.T) {
	text := tgbotapi.NewMessage(1092208354, "Skip this message!")
	_, err := testTgBotAPI.Send(text)
	require.NoError(t, err)
}
