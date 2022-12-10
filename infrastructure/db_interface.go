package infrastructure

import (
	"era-light-bot/domain"
)

type Db interface {
	GetUser(id int64) (*domain.User, error)
	GetUserIds() ([]int64, error)
	SaveUser(user domain.User) error
	GetLightState() (*domain.LightState, error)
	SaveLightState(state *domain.LightState) error
}
