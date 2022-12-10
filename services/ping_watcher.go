package services

import (
	"fmt"
	"log"
	"tg-light-bot/common"
	"tg-light-bot/domain"
	"tg-light-bot/infrastructure"
	"time"
)

type PingChannel chan bool

func NewPingChannel() *PingChannel {
	c := make(PingChannel)
	return &c
}

type PingWatcher struct {
	c            *PingChannel
	ticker       *time.Ticker
	conn         infrastructure.Db
	currentState *domain.LightState
	lastPing     time.Time
	bot          *EraBot
	config       common.PingConfig
}

func NewPingWatcher(pingChannel *PingChannel, bot *EraBot, dbConnection infrastructure.Db, config common.PingConfig) *PingWatcher {
	currentState, err := dbConnection.GetLightState()

	log.Printf("[ping_watcher] Init pingWatcher with last state %v\n", currentState)

	if err != nil {
		log.Fatalf("[ping_watcher] Error loading last light state. Message: %s\n", err)
	}

	return &PingWatcher{
		c:            pingChannel,
		ticker:       time.NewTicker(time.Second),
		conn:         dbConnection,
		currentState: currentState,
		lastPing:     common.GetNow(),
		bot:          bot,
		config:       config,
	}
}

func (p *PingWatcher) Watch() {
	go func() {
		for {
			select {
			case <-*p.c:
				p.lastPing = common.GetNow()
				if p.currentState.IsOn == false { // so it was OFF and is a first ping meaning that is ON now
					oldDate := p.currentState.ChangedAt
					p.switchState()
					p.bot.MassSend(fmt.Sprintf(
						"Дали світло!\nМи були без світла %s\nПеріод: %s - %s",
						common.FormatDuration(p.currentState.ChangedAt.Sub(oldDate)),
						oldDate.Format("15:04 02.01.2006"),
						p.currentState.ChangedAt.Format("15:04 02.01.2006"),
					))
				}
			case <-p.ticker.C:
				if p.currentState.IsOn && common.GetNow().Sub(p.lastPing) > p.config.Threshold {
					oldDate := p.currentState.ChangedAt
					p.switchState()
					p.bot.MassSend(fmt.Sprintf(
						"Світло відсутнє.\nМи були зі світлом %s\nПеріод: %s - %s",
						common.FormatDuration(p.currentState.ChangedAt.Sub(oldDate)),
						oldDate.Format("15:04 02.01.2006"),
						p.currentState.ChangedAt.Format("15:04 02.01.2006"),
					))
				}
			}
		}
	}()
}

func (p *PingWatcher) switchState() *domain.LightState {
	newState := &domain.LightState{IsOn: !p.currentState.IsOn, ChangedAt: common.GetNow()}
	log.Println("New State: ", newState)
	err := p.conn.SaveLightState(newState)
	if err != nil {
		log.Fatalf("[ping_watcher] Error saving last light state. Message: %s\n", err)
	}
	p.currentState = newState
	return newState
}
