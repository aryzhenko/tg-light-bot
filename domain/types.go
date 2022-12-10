package domain

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

type User struct {
	tgbotapi.User
	CreatedAt    time.Time
	LastActivity time.Time
	IsNew        bool
}

type LightState struct {
	IsOn      bool
	ChangedAt time.Time
}
