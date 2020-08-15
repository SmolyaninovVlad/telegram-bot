package bot

import( 
	"fmt"
	"strconv"
	"telegaBot/structers"
	"telegaBot/post"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

const tgbotapiKey = "901228271:AAGEfVEGO4TXI8sZfmKQPMDYIfNQxPz-Qxs"

var (
	mainMenu = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("ℹ Информация"),
			tgbotapi.NewKeyboardButton("🗒 Создать заявку"),
		),
	)
	requestMenu = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("ℹ Информация(запрос)"),
			tgbotapi.NewKeyboardButton("❌ Отмена"),
		),
	)
)

var (
	newRequestMap = make(map[int64]*Structers.NewRequest)
	customFieldFirst = new(Structers.Customf)
	customFieldSecond = new(Structers.Customf)
	bot *tgbotapi.BotAPI
	err error
	updChannel tgbotapi.UpdatesChannel
	update tgbotapi.Update
	updConfig tgbotapi.UpdateConfig
	botUser tgbotapi.User
)

func InitBot() tgbotapi.UpdatesChannel{
	bot, err = tgbotapi.NewBotAPI(tgbotapiKey)
	if err != nil {
		panic("bot init error:" + err.Error())
	}
	botUser, err = bot.GetMe()
	if err != nil {
		panic("bot getMe error:" + err.Error())
	}
	fmt.Printf("auth is ok! Bot is : %s\n", botUser.FirstName)
	updConfig.Timeout = 60
	updConfig.Limit = 1
	updConfig.Offset = 0

	updChannel, err = bot.GetUpdatesChan(updConfig)
	if err != nil {
		panic("updChannel error:" + err.Error())
	}
	return updChannel
}
func InitBotCycle(updChannel tgbotapi.UpdatesChannel) error{
	// Начальные константные значения редмайна
	customFieldFirst.Id = "3"
	customFieldFirst.Name = "Подразделение-инициатор"
	customFieldFirst.Value = "Внешний канал"
	customFieldSecond = new(Structers.Customf)
	customFieldSecond.Id = "6"
	customFieldSecond.Name = "ExternalChannelUserId"
	//Цикл для обработки сообщений
	for {
		update = <-updChannel
		if update.Message != nil {
			//Если команда
			if update.Message.IsCommand(){
				cmdText := update.Message.Command()
				if cmdText == "help" {
					msgConfig := tgbotapi.NewMessage(
						update.Message.Chat.ID,
						"Чтобы принять запрос на сопровождение введите команду /menu и выберите интересующий вас пункт")
					bot.Send(msgConfig)
				} else if cmdText == "menu" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Главное меню")
					msg.ReplyMarkup = mainMenu
					bot.Send(msg)
				} else if cmdText == "start" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать. Чем я могу вам помочь:\n/menu - Главное меню\n/help - Дополнительная информация\nВедутся работы...")
					bot.Send(msg)
				} 
			} else {
				switch update.Message.Text {
				case mainMenu.Keyboard[0][0].Text: //Инфо (Главная)
					infoMain()
				case mainMenu.Keyboard[0][1].Text: //Новый (Запрос)
					newRequest()
				case requestMenu.Keyboard[0][0].Text: // Инфо (Запрос)
					infoRequest()
				case requestMenu.Keyboard[0][1].Text: // Отмена (Запрос)
					delete(newRequestMap, update.Message.Chat.ID)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Операция отменена")
					msg.ReplyMarkup = mainMenu
					bot.Send(msg)
				default:
					requestProcessing()
				}
			}
		}

	}
}
func SendMessageToUser(chat_id int64, message string) {
	msgConfig := tgbotapi.NewMessage(chat_id, message)
	bot.Send(msgConfig)
}
func infoMain(){
	message := "Я бот ПримСоцБанка. Я принимаю запроса на сопровождение от пользователей телеграма.\n В порядке очереди наши специалисты обработают ваш запрос и дадут обратную связь. Я помогу вам создать запрос и сообщу как только он будет выполнен."
	msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	bot.Send(msgConfig)
}
func infoRequest(){
	message := "Для создания запроса вам необходимо заполнить тему запроса и описание. В кротчайшие сроки специалисты обработают ваш запрос. По завершению работы или если будет дополнительная информация от наших специалистов, я вам сообщу\n"
	msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	bot.Send(msgConfig)
}
func sendQuestion(state int) string {
	var result string
	switch state {
		case 0: result = "Введите тему: "
		case 1: result = "Введите Описание: "
		case -1: result = "⏳ Ожидайте, идёт обработка запроса"
		case -2: result = "✅ Ваш запрос принят\nКогда ваш запрос будет готов или если будут какие-то комментарии - я вам сообщу, ожидайте"
		case -3: result = "🛠 Извините, сервис временно недоступен."
	}
	return result
}
func newRequest() {
	newRequestMap[update.Message.Chat.ID] = new(Structers.NewRequest)
	//Начальное заполнение константных полей
	newRequestMap[update.Message.Chat.ID].State = 0
	//Добавляем в переменную запроса данные о кастомных полях
	var array []*Structers.Customf
	//переменная о внешнем канале
	array = append(array, customFieldFirst)
	//переменная об айди юзера
	customFieldSecond.Value = strconv.FormatInt(update.Message.Chat.ID, 10)
	array = append(array, customFieldSecond)
	newRequestMap[update.Message.Chat.ID].Custom_fields = array
	newRequestMap[update.Message.Chat.ID].Project_id = "16"
	newRequestMap[update.Message.Chat.ID].Tracker_id = "2"
	newRequestMap[update.Message.Chat.ID].Status_id = "1"
	newRequestMap[update.Message.Chat.ID].Priority_id = "2"
	msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, sendQuestion(0))
	msgConfig.ReplyMarkup = requestMenu
	// msgConfig. ДОБАВИТЬ УДАЛЕНИЕ МЕНЮ
	bot.Send(msgConfig)
}
func requestProcessing(){
	//Получаем текущий опрос и обрабатываем
	cs, ok := newRequestMap[update.Message.Chat.ID]
	if ok {
		switch cs.State {
			//Заполнили первый, переход на второй шаг
			case 0: {
				cs.State = 1
				cs.Subject = update.Message.Text
				msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, sendQuestion(cs.State))
				bot.Send(msgConfig)
			}
			//Заполнили второй, переход на третий шаг (конец)
			case 1: {
				cs.State = -1
				cs.Description = update.Message.Text
				msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, sendQuestion(cs.State))
				// отправляем запрос на создание
				bot.Send(msgConfig)
				err = post.SendPost(cs)
				if err != nil {
					//Ошибка отправки, сообьщение пользователю об ошибке
					fmt.Printf("send post error: %v\n", err)
					msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, sendQuestion(-3))
					msgConfig.ReplyMarkup = mainMenu
					bot.Send(msgConfig)
				} else {
					msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, sendQuestion(-2))
					msgConfig.ReplyMarkup = mainMenu
					bot.Send(msgConfig)
				}
				//Удаление опроса
				delete(newRequestMap, update.Message.Chat.ID)
			}
		}
	} else {
		//Other messages
		// логирование сообщений которые не попали в обработчик
		fmt.Printf("form: %s; message: %s\n", update.Message.From.FirstName, update.Message.Text)
	}
}
