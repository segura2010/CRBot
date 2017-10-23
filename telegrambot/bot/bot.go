package bot


import (
    "time"
    "fmt"
    "strings"
    //"log"
    "strconv"
    "encoding/hex"
    
    "github.com/tucnak/telebot"

    "github.com/segura2010/cr-go/packets"
    "github.com/segura2010/cr-go/resources"
    "github.com/segura2010/cr-go/utils"

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
        // log.Printf("Received message..")
        userId := fmt.Sprintf("%d", message.Sender.ID)
        chatId := fmt.Sprintf("%d", message.Chat.ID)

        if strings.Index(message.Text, "/get ") == 0 {
            playerTag := clearTag(message.Text[5:])
            if utils.IsValidTag(playerTag){
                db.AddPlayerStatsJob(playerTag, chatId)
                myBot.Bot.SendMessage(message.Chat, "Retrieving stats...", nil)
            }else{
                myBot.Bot.SendMessage(message.Chat, "Invalid player tag", nil)
            }
        }else if strings.Index(message.Text, "/save ") == 0 {
            playerTag := clearTag(message.Text[6:])
            if utils.IsValidTag(playerTag){
                err := db.AddPlayerTag(playerTag, userId)
                if err != nil{
                    myBot.Bot.SendMessage(message.Chat, "There was an error saving your user tag :(", nil)
                    continue
                }
                myBot.Bot.SendMessage(message.Chat, "Your player tag was saved!", nil)
            }else{
                myBot.Bot.SendMessage(message.Chat, "Invalid player tag", nil)
            }
        }else if strings.Index(message.Text, "/delete") == 0 {
            err := db.DeletePlayerTag(userId)
            if err != nil{
                myBot.Bot.SendMessage(message.Chat, "There was an error deleting your user tag :(", nil)
                continue
            }
            myBot.Bot.SendMessage(message.Chat, "Your player tag was deleted!", nil)
        }else if strings.Index(message.Text, "/me") == 0 {
            playerTag := db.GetPlayerTag(userId)
            if playerTag == ""{
                myBot.Bot.SendMessage(message.Chat, "There was an error retrieving your user tag :(", nil)
                continue
            }
            db.AddPlayerStatsJob(playerTag, chatId)
            myBot.Bot.SendMessage(message.Chat, "Retrieving stats...", nil)
        }else{
            // help..
            r := fmt.Sprintf("Available commands:\n\t/get PLAYER_TAG : get stats for the specified player")
            r += fmt.Sprintf("\n\t/save PLAYER_TAG : saves your player tag")
            r += fmt.Sprintf("\n\t/delete : delete your saved player tag")
            r += fmt.Sprintf("\n\t/me : get your stats (based on the saved tag)")
            myBot.Bot.SendMessage(message.Chat, r, nil)
        }
    }
}

func clearTag(tag string) (string){
    tag = strings.Replace(tag, ":", "", -1)
    tag = strings.Replace(tag, " ", "", -1)
    tag = strings.Replace(tag, "#", "", -1)
    tag = strings.ToUpper(tag)

    return tag
}

func formatChestsOrder(playerInfo packets.ServerVisitHome)(string){
    result := fmt.Sprintf("*Next chests*:\n|")

    for i:=0;i<8;i++{
        chestPos := (playerInfo.ChestCycle.CurrentPosition + int32(i)) % int32(len(resources.ChestOrder))
        chest := resources.ChestOrder[chestPos]
        result += fmt.Sprintf("%s|", chest)
    }

    // get next Giant and Magic chest
    i := 0
    magicPos := int32(-1)
    giantPos := int32(-1)
    for{
        chestPos := (playerInfo.ChestCycle.CurrentPosition + int32(i)) % int32(len(resources.ChestOrder))
        chest := resources.ChestOrder[chestPos]
        
        if chest == "Magic"{
            magicPos = (playerInfo.ChestCycle.CurrentPosition + int32(i)) - playerInfo.ChestCycle.CurrentPosition
        }else if chest == "Giant"{
            giantPos = (playerInfo.ChestCycle.CurrentPosition + int32(i)) - playerInfo.ChestCycle.CurrentPosition
        }

        if magicPos > -1 && giantPos > -1{
            break
        }

        i += 1
    }

    result += fmt.Sprintf("\nNext *Magical* in *%d* wins", magicPos)
    result += fmt.Sprintf("\nNext *Giant* in *%d* wins", giantPos)

    return result
}

func formatUserStats(playerInfo packets.ServerVisitHome) (string){
    var result string
    winsPlusLosses := float32(playerInfo.Wins + playerInfo.Losses)

    result = fmt.Sprintf("*%s* 🏆%d (Record: %d)", playerInfo.Username, playerInfo.Trophies, playerInfo.Stats.RecordTrophies)
    
    if playerInfo.HasClan{
        result += fmt.Sprintf("\n🛡 %s", playerInfo.Clan.Name)
    }
    
    result += fmt.Sprintf("\n*Level*: %d | *Cards Found*: %d", playerInfo.Level, playerInfo.Stats.UnlockedCards)
    result += fmt.Sprintf("\n💰 %d | 💎 %d", playerInfo.Gold, playerInfo.Gems)
    result += fmt.Sprintf("\n*Wins*: %d | *Losses*: %d", playerInfo.Wins, playerInfo.Losses)
    result += fmt.Sprintf("\n*Win Rate* %.2f%% | *Games*: %d", (float32(playerInfo.Wins) / winsPlusLosses)*100.0, playerInfo.Games )
    result += fmt.Sprintf("\n*3 👑 Wins*: %d (%.2f%%) | *Donations*: %d", playerInfo.Stats.CrownWins3, (float32(playerInfo.Stats.CrownWins3) / float32(playerInfo.Wins))*100.0, playerInfo.Stats.Donations)
    result += fmt.Sprintf("\n*Challenge Cards Won*: %d | *Challenge Max Wins*: %d", playerInfo.Stats.ChallengeCardsWon, playerInfo.Stats.ChallengeMaxWins)
    result += fmt.Sprintf("\n*Tournament Games*: %d", playerInfo.TournamentGames)
    result += fmt.Sprintf("\n--- Chests ---")
    result += fmt.Sprintf("\n%s", formatChestsOrder(playerInfo))
    result += fmt.Sprintf("\nNext *SuperMagical* in *%d* wins", (playerInfo.ChestCycle.SuperMagicalPos-playerInfo.ChestCycle.CurrentPosition))
    result += fmt.Sprintf("\nNext *Legendary* in *%d* wins", (playerInfo.ChestCycle.LegendaryPos-playerInfo.ChestCycle.CurrentPosition))
    result += fmt.Sprintf("\nNext *Epic* in *%d* wins", (playerInfo.ChestCycle.EpicPos-playerInfo.ChestCycle.CurrentPosition))
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

        chatId, err := strconv.ParseInt(jobInfo[1], 10, 64)
        if err != nil{
            continue
        }

        if jobInfo[2] == "err"{
            SendMessage(chatId, "There was an error retreiving the stats for that player :(", nil)
            continue
        }

        tobytes, err := hex.DecodeString(jobInfo[2])
        if err != nil{
            continue
        }
        playerInfo := packets.NewServerVisitHomeFromBytes(tobytes)
        // log.Printf("Received stats.. %v", playerInfo)
        
        msg := formatUserStats(playerInfo)
        sendOptions := telebot.SendOptions{
            ParseMode: "Markdown",
        }
        SendMessage(chatId, msg, &sendOptions)
    }
}


