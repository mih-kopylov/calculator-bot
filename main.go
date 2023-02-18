package main

import (
	"fmt"
	"github.com/caarlos0/env/v7"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joomcode/errorx"
	"github.com/mnogu/go-calculator"
	"github.com/sirupsen/logrus"
)

type Config struct {
	ApiToken string `env:"TELEGRAM_API_TOKEN,required"`
}

type Status struct {
	examplesCount int
}

func main() {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		logrus.Fatal(errorx.EnhanceStackTrace(err, "failed to parse config"))
	}
	status := &Status{}

	bot, err := tgbotapi.NewBotAPI(cfg.ApiToken)
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		m := update.Message
		if m == nil {
			m = update.ChannelPost
		}
		if m == nil {
			logrus.Debug("update without message received")
			continue
		}

		responseText := prepareResponse(m, status)

		message := tgbotapi.NewMessage(m.Chat.ID, responseText)
		message.ReplyToMessageID = m.MessageID

		_, err = bot.Send(message)
		if err != nil {
			logrus.Error(errorx.EnhanceStackTrace(err, "failed to send a reply"))
		}
	}
}

func prepareResponse(message *tgbotapi.Message, status *Status) string {
	if message.IsCommand() {
		switch message.Command() {
		case "count":
			return fmt.Sprintf("Вычислено примеров : %v", status.examplesCount)
		default:
			return `Доступные команды:
/help - показывает эту справку
/count - показывает количество выполненных примеров`
		}
	}
	result, err := calculator.Calculate(message.Text)
	if err != nil {
		return "Невозможно вычислить результат"
	} else {
		status.examplesCount++
		return fmt.Sprintf("%g", result)
	}

}
