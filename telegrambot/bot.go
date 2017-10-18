package bot


import (
    "time"
    "fmt"
    
    "github.com/tucnak/telebot"

    "CRBot/db"
)

type TelegramBot struct{
	Bot *telebot.Bot
    Token string
    Started bool
}

var instance *TelegramBot = nil

func CreateInstance(token, name string) *TelegramBot {
    instance = &TelegramBot{Token:token, Name:name, Started:false}
    bot, err := telebot.NewBot(token)
    if err != nil {
        return nil
    }

    instance.Started = true
    instance.Bot = bot
    go listenMessages()

    return instance
}

func GetInstance() *TelegramBot {
    return instance
}

func RefreshSession(){
    CreateInstance(instance.Token, instance.Name)
}


func listenMessages(){
    myBot := GetInstance()
    messages := make(chan telebot.Message)
    myBot.Bot.Listen(messages, 1*time.Second)

    for message := range messages {
        userID := fmt.Sprintf("%d", message.Sender.ID)
        
        if message.Text == "/get" {
            response := fmt.Sprintf("Stats: ...")
            myBot.Bot.SendMessage(message.Chat, response, nil)
        }else{
            // help..
            r := fmt.Sprintf("Help:...")
            myBot.Bot.SendMessage(message.Chat, r, nil)
        }
    }
}

func SendMessage(to int64, message string){
    myBot := GetInstance()

    chat := telebot.Chat{ID: to}
    myBot.Bot.SendMessage(chat, message, nil)
}
