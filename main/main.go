package main

import( 
	"telegaBot/structers"
	"telegaBot/incomingApi"
	"telegaBot/bot"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	customFieldFirst, customFieldSecond *Structers.Customf
	updChannel tgbotapi.UpdatesChannel
)
func main() {
	//инициализация API в отдельной горутине
	go func () {
		api.MainAPI()
	}()
	updChannel = bot.InitBot()
	bot.InitBotCycle(updChannel)

}

