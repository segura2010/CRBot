package bot


import (
    "time"
    "fmt"
    "strings"
    "log"
    "strconv"
    "encoding/hex"
    
    "github.com/tucnak/telebot"

    "github.com/segura2010/cr-go/packets"

    "CRBot/db"
)

type TelegramBot struct{
	Bot *telebot.Bot
    Token string
    Started bool
}

var instance *TelegramBot = nil

func CreateInstance(token string) *TelegramBot {
    instance = &TelegramBot{Token:token, Started:false}
    bot, err := telebot.NewBot(token)
    if err != nil {
        panic(err)
    }

    instance.Started = true
    instance.Bot = bot
    go listenMessages()
    go listenStats() // suscribe to redis channel to receive and send the stats results

    return instance
}

func GetInstance() *TelegramBot {
    return instance
}

func RefreshSession(){
    CreateInstance(instance.Token)
}

func SendMessage(to int64, message string, options *telebot.SendOptions){
    myBot := GetInstance()

    chat := telebot.Chat{ID: to}
    myBot.Bot.SendMessage(chat, message, options)
}

func listenMessages(){
    myBot := GetInstance()
    messages := make(chan telebot.Message)
    myBot.Bot.Listen(messages, 1*time.Second)

    for message := range messages {
        log.Printf("Received message..")
        //userID := fmt.Sprintf("%d", message.Sender.ID)
        
        if strings.Index(message.Text, "/get ") == 0 {
            playerTag := strings.Replace(message.Text[5:], ":", "", -1)

            chatId := fmt.Sprintf("%d", message.Chat.ID)
            db.AddPlayerStatsJob(playerTag, chatId)

            response := fmt.Sprintf("Stats: ...")
            myBot.Bot.SendMessage(message.Chat, response, nil)
        }else{
            // help..
            r := fmt.Sprintf("Available commands:\n/get PLAYER_TAG")
            myBot.Bot.SendMessage(message.Chat, r, nil)
        }
    }
}

func formatUserStats(playerInfo packets.ServerVisitHome) (string){
    var result string
    result = fmt.Sprintf("*%s* 🏆%d (Record: %d)", playerInfo.Username, playerInfo.Trophies, playerInfo.Stats.RecordTrophies)
    result += fmt.Sprintf("\n*Clan*: %s", playerInfo.Clan.Name)
    result += fmt.Sprintf("\n*Level*: %d | *Cards Found*: %d", playerInfo.Level, playerInfo.Stats.UnlockedCards)
    result += fmt.Sprintf("\n*Gold*: %d | *Gems*: %d", playerInfo.Gold, playerInfo.Gems)
    result += fmt.Sprintf("\n*Wins*: %d | *Losses*: %d | *Games*: %d", playerInfo.Wins, playerInfo.Losses, playerInfo.Games)
    result += fmt.Sprintf("\n*3 Crowns Wins*: %d | *Donations*: %d", playerInfo.Stats.CrownWins3, playerInfo.Stats.Donations)
    result += fmt.Sprintf("\n*Challenge Cards Won*: %d | *Challenge Max Wins*: %d", playerInfo.Stats.ChallengeMaxWins, playerInfo.Stats.ChallengeCardsWon)
    result += fmt.Sprintf("\n*Tournament Games*: %d", playerInfo.TournamentGames)
    result += fmt.Sprintf("\n--- Chests ---")
    result += fmt.Sprintf("\nNext *SuperMagical* in *%d* wins", (playerInfo.ChestCycle.SuperMagicalPos-playerInfo.ChestCycle.CurrentPosition))
    result += fmt.Sprintf("\nNext *Legendary* in *%d* wins", (playerInfo.ChestCycle.LegendaryPos-playerInfo.ChestCycle.CurrentPosition))
    //result += fmt.Sprintf("\nNext *Epic* in *%d* days", (playerInfo.ChestCycle.MagicalPos-playerInfo.ChestCycle.CurrentPosition))
    result += fmt.Sprintf("\n--- Shop ---")
    result += fmt.Sprintf("\nNext *Legendary* in *%d* days", (playerInfo.ShopOffers.Legendary-playerInfo.ShopOffers.CurrentDay))
    result += fmt.Sprintf("\nNext *Epic* in *%d* days", (playerInfo.ShopOffers.Epic-playerInfo.ShopOffers.CurrentDay))
    result += fmt.Sprintf("\nNext *Arena* in *%d* days", (playerInfo.ShopOffers.Arena-playerInfo.ShopOffers.CurrentDay))

    return result
}

func listenStats(){
    //myBot := GetInstance()
    pubsub := db.SubscribeToPlayerStats()
    defer pubsub.Close()
    messages := pubsub.Channel()

    for message := range messages {
        jobInfo := strings.Split(message.Payload, ":")
        tobytes, err := hex.DecodeString(jobInfo[2])
        if err != nil{
            continue
        }
        playerInfo := packets.NewServerVisitHomeFromBytes(tobytes)
        //log.Printf("Received stats.. %v", playerInfo)
        chatId, err := strconv.ParseInt(jobInfo[1], 10, 64)
        if err != nil{
            continue
        }
        msg := formatUserStats(playerInfo)
        sendOptions := telebot.SendOptions{
            ParseMode: "Markdown",
        }
        SendMessage(chatId, msg, &sendOptions)
    }
}


