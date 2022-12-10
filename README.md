# What is tg-light-bot
tg-light-bot is a bot to inform users if light is on or of at the moment giving additional information about time whe it was on/of and total duration of current state.
Also, it notifies when state is changed giving information about how long it was in last state.

## How it works
Bot must have permanent IP address (or DNS name), so agent can make requests to it.
Anything able to make HTTP requests can be an agent. You can use your home computer, router, or any hardware able to periodically send HTTP request. I use my router sending requests with wget by cron.

## Requirements
You need to install MySQL server with structures created (can be found in docker-compose/mysql/init/db.sql).

## Install and run
1. Instal dependencies by running ```make up```
2. Build code by running ```make build```
3. Upload binary ```build/light-bot``` to your remote server
4. Copy ```.env.dist``` to ```.env``` and set configuration values
5. Upload ```.env``` to the same folder as binary
6. Run binary
7. MAke your agent send request GET http://BOT_IP:LIGHT_BOT_PING_PORT/ping periodically (i.e. every minute by cron)

# Configuration
You can configure bot via .env file or environment variables. List of available options in table below

| Option name              | Type   | Description                                                                                                                                 |
|--------------------------|--------|---------------------------------------------------------------------------------------------------------------------------------------------|
| LIGHT_BOT_TOKEN          | String | Telegram bot token. You can obtain it using special telegram bot [Botfather](https://t.me/BotFather)                                        |
| LIGHT_BOT_DB_HOST        | String | MySQL database host. IP address ot DNS name                                                                                                 |
| LIGHT_BOT_DB_USER        | String | MySQL database user                                                                                                                         |
| LIGHT_BOT_DB_PASSWORD    | String | MySQL database password                                                                                                                     |
| LIGHT_BOT_DB_NAME        | String | MySQL database name                                                                                                                         |
| LIGHT_BOT_PING_PORT      | Int    | Port that bot will listen for incoming agent requests                                                                                       |
| LIGHT_BOT_PING_THRESHOLD | String | Period between last ping and considering light is gone. **MUST NOT** be less then period of agent sending requests. Value examples: 10s, 5m |

