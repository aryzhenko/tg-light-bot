package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"tg-light-bot/common"
	"tg-light-bot/infrastructure/mysql"
	"tg-light-bot/services"
)

var PingChan *services.PingChannel

func main() {

	config, err := common.LoadConfig(".env")

	if err != nil {
		log.Fatalf("Unable to load config. %s\n", err)
	}

	conn := mysql.NewConnection(config.Db)

	PingChan = services.NewPingChannel()
	router := gin.Default()

	router.GET("/ping", Ping)

	bot := services.NewEraBot(
		config.BotToken,
		conn,
		services.NewLight(conn),
	)
	bot.Start()

	pw := services.NewPingWatcher(PingChan, bot, conn, config.Ping)
	pw.Watch()

	router.Run(fmt.Sprintf(":%d", config.Ping.Port))
}

// TODO move it somewhere
func Ping(c *gin.Context) {
	*PingChan <- true
	c.IndentedJSON(http.StatusOK, "pong")
}
