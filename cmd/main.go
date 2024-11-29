package main

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/bolotskihUlyanaa/christmasBot/pkg"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/spf13/viper"
)

type Config struct {
	token        string
	personFile   string
	logFile      string
	startMessage string
	rulesMessage string
	startStiker  string
}

func loadPersonsFromFile() {
	f, err := os.Open(conf.personFile)
	if err != nil {
		log.Panicln("Error with open file: ", err.Error())
	}
	data, err := io.ReadAll(f)
	f.Close()
	if err != nil {
		log.Panicln("Error with read file: ", err.Error())
	}
	err = json.Unmarshal(data, &players)
	if err != nil {
		log.Panicln("Error with decoding: ", err.Error())
	}
}

func savePersonsToFile() {
	f, err := os.Create(conf.personFile)
	if err != nil {
		log.Panicln("Error with open file: ", err.Error())
	}
	b, err := json.Marshal(players)
	if err != nil {
		log.Panicln("Error with encoding: ", err.Error())
	}
	_, err = f.Write(b)
	f.Close()
	if err != nil {
		log.Panicln("Error with write: ", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

var conf Config
var players pkg.Persons

func main() {
	//инициализаируем данные из конфигурационного файла
	if err := initConfig(); err != nil {
		log.Fatal("Error initializing configs: ", err.Error())
	}
	conf = Config{
		token:        viper.GetString("token"),
		personFile:   viper.GetString("personFile"),
		logFile:      viper.GetString("logFile"),
		startMessage: viper.GetString("startMessage"),
		rulesMessage: viper.GetString("rulesMessage"),
		startStiker:  viper.GetString("startStiker"),
	}

	//O_CREATE - создание нового файла если не существует
	//O_WRONLY - только для записи
	//O_APPEND - добавление данных при записи
	//FileMode - 0xxx
	//0 - восьмеричная система
	//x:
	//4 - чтение
	//2 - запись
	//1 - выполнение
	//0 - полный запрет
	file, err := os.OpenFile(conf.logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error failed to open log file:", err)
	}
	log.SetOutput(file)
	log.Println(time.Now())
	loadPersonsFromFile()
	count := len(players.Persons)

	bot, err := telego.NewBot(conf.token, telego.WithDefaultDebugLogger())
	if err != nil {
		log.Panicln("Error create bot: ", err)
	}

	//создаем 2 кнопки
	keyboard := tu.Keyboard(tu.KeyboardRow(
		tu.KeyboardButton("Узнать получателя"),
		tu.KeyboardButton("Узнать правила")),
	).WithResizeKeyboard()

	updates, _ := bot.UpdatesViaLongPolling(nil)
	bh, _ := th.NewBotHandler(bot, updates)
	defer bh.Stop()
	defer bot.StopLongPolling()

	//отправляем приветственное письмо, реагирует на /start
	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		chatID := tu.ID(update.Message.Chat.ID)
		_, _ = bot.SendMessage(tu.Message(
			chatID,
			conf.startMessage,
		))
		_, _ = bot.SendSticker(tu.Sticker(
			chatID,
			tu.FileFromID(conf.startStiker)).WithReplyMarkup(keyboard),
		)
	}, th.CommandEqual("start"))

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		chatID := tu.ID(update.Message.Chat.ID)

		//определяем по нику пользователя
		person := players.Find(update.Message.From.Username)

		//если друг еще не назначен
		//если чел повторно отправляет, то ему тот же чел попадается
		if person.Friend == "" {
			var idx int
			for {
				idx = rand.Intn(count)
				if person.Filter2(&players.Persons[idx]) {
					savePersonsToFile()
					break
				}
			}
		}

		_, _ = bot.SendMessage(
			tu.Message(
				chatID,
				"<tg-spoiler>@"+person.Friend+"</tg-spoiler>",
			).WithParseMode("HTML").WithReplyMarkup(keyboard),
		)
	}, th.TextEqual("Узнать получателя"))

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		chatID := tu.ID(update.Message.Chat.ID)
		_, _ = bot.SendMessage(
			tu.Message(
				chatID,
				conf.rulesMessage,
			).WithParseMode("HTML").WithReplyMarkup(keyboard),
		)
	}, th.TextEqual("Узнать правила"))
	bh.Start()
}
