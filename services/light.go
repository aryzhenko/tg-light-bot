package services

import (
	"fmt"
	"tg-light-bot/common"
	"tg-light-bot/infrastructure"
)

type Light struct {
	db infrastructure.Db
}

func NewLight(db infrastructure.Db) *Light {
	return &Light{
		db: db,
	}
}

func (l *Light) GetState() (string, error) {
	state, err := l.db.GetLightState()
	if err != nil {
		return "", err
	}
	diff := common.FormatDuration(common.GetNow().Sub(state.ChangedAt))
	if state.IsOn {
		return fmt.Sprintf("Світло є\nУвімкнули об %s (%s тому)", state.ChangedAt.Format("15:04 02.01.2006"), diff), nil
	} else {
		return fmt.Sprintf("Світла немає\nВимкнули об %s (%s тому)", state.ChangedAt.Format("15:04 02.01.2006"), diff), nil
	}
}
