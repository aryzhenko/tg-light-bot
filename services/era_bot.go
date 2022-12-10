package services

import (
	"era-light-bot/common"
	"era-light-bot/domain"
	"era-light-bot/infrastructure"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

const cmdLight = "Чи є світло?"

var lightKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(cmdLight),
	),
)

type EraBot struct {
	bot          *tgbotapi.BotAPI
	lightService *Light
	db           infrastructure.Db
}

func NewEraBot(token string, db infrastructure.Db, lightService *Light) *EraBot {
	bot, err := tgbotapi.NewBotAPI(token)
	if nil != err {
		log.Println("[era_bot] Error initializing bot. Message: ", err)
	}

	return &EraBot{
		bot:          bot,
		db:           db,
		lightService: lightService,
	}

}

func (e *EraBot) Start() {
	go func() {
		updateConfig := tgbotapi.NewUpdate(0)
		updateConfig.Timeout = 30

		updates := e.bot.GetUpdatesChan(updateConfig)

		for update := range updates {
			if update.Message == nil {
				continue
			}

			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "start":
					e.handleStartCommand(update.Message)
				}
			}

			switch update.Message.Text {
			case cmdLight:
				response := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				response.ReplyMarkup = lightKeyboard
				if text, err := e.lightService.GetState(); err != nil {
					response.Text = "Невідома помилка"
					log.Println("[era_bot] Error loading last light state. ", err)
				} else {
					response.Text = text
				}

				if _, err := e.bot.Send(response); err != nil {
					log.Println("[era_bot] Error sending message. ", err)
				}
			}
		}
	}()
}
func (e *EraBot) MassSend(text string) {
	userIds, err := e.db.GetUserIds()
	if err != nil {
		log.Printf("[era_bot] Error loading user ids: %s\n", err)
		return
	}
	for _, id := range userIds {
		message := tgbotapi.NewMessage(id, text)
		_, err := e.bot.Send(message)

		if err != nil {
			log.Printf("[era_bot] Error sending message to user [%d]: %s\n", id, err)
		}
	}
}

func (e *EraBot) handleStartCommand(msg *tgbotapi.Message) {
	responseText := fmt.Sprintf("Тепер я буду повідомляти Вас коли вимкнуть або увімкнуть світло.\n" +
		"Також Ви зможете дізнатися чи є зараз світло за допомогою кнопки нижче\n" +
		"Приємного користування")
	response := tgbotapi.NewMessage(msg.Chat.ID, responseText)
	response.ReplyMarkup = lightKeyboard
	user := domain.User{
		User:         *msg.From,
		CreatedAt:    common.GetNow(),
		LastActivity: common.GetNow(),
		IsNew:        true,
	}
	err := e.db.SaveUser(user)
	if err != nil {
		log.Printf("[era_bot] Error saving user. %s", err)
	}

	if _, err := e.bot.Send(response); err != nil {
		log.Println("[era_bot] Error sending message. ", err)
	}

}
