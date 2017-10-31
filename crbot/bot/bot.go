package bot


import (
    "log"
    "encoding/hex"
    "errors"
    
    "github.com/tucnak/telebot"

    "github.com/segura2010/cr-go/client"
    "github.com/segura2010/cr-go/packets"
    "github.com/segura2010/cr-go/utils"

    //"CRBot/db"
    "CRBot/config"
)

type TelegramBot struct{
	Bot *telebot.Bot
    Token string
    Started bool
}

func GetStatsForPlayer(tag string) (packets.Packet, error){
    var err error
    serverKey, err := hex.DecodeString(config.GetInstance().CRBot.ServerKey)
    if err != nil{
        log.Printf("Invalid server key")
    }

    helloPayload := packets.NewDefaultClientHello()
    helloPayload.ContentHash = config.GetInstance().CRBot.ContentHash
    helloPkt := packets.Packet{
        Type: packets.MessageType["ClientHello"],
        Version: 0,
        Payload: helloPayload.Bytes(),
    }

    loginPayload := packets.NewDefaultClientLogin()
    HiLo := utils.Tag2HiLo(config.GetInstance().CRBot.PlayerTag) // test account
    loginPayload.Hi = HiLo[0]
    loginPayload.Lo = HiLo[1]
    loginPayload.PassToken = config.GetInstance().CRBot.PassToken
    loginPayload.ContentHash = config.GetInstance().CRBot.ContentHash

    loginPkt := packets.Packet{
        Type: packets.MessageType["ClientLogin"],
        Version: 0,
        DecryptedPayload: loginPayload.Bytes(),
    }

    var basicPkt packets.Packet

    serverAddress := "game.clashroyaleapp.com:9339"

    c := client.NewCRClient(serverKey)
    c.Connect(serverAddress)
    defer c.Close()

    c.SendPacket(helloPkt)
    basicPkt = c.RecvPacket() // receive hello response
    if basicPkt.Type == packets.MessageType["ServerLoginFailed"]{
        return packets.Packet{}, errors.New("ServerLoginFailed on Hello")
    }

    c.SendPacket(loginPkt)
    basicPkt = c.RecvPacket() // receive login response
    if basicPkt.Type == packets.MessageType["ServerLoginFailed"]{
        return packets.Packet{}, errors.New("ServerLoginFailed on Login")
    }
    //loginOk := packets.NewServerLoginOkFromBytes(basicPkt.DecryptedPayload)
    log.Printf("LoginOk")

    // receive multiple packets the server sends after login
    c.RecvPacket() // OwnHomeData
    c.RecvPacket() // InboxGlobal
    c.RecvPacket() // FriendList

    // send keepalive
    basicPkt = packets.Packet{
        Type: packets.MessageType["ClientKeepAlive"],
        Version: 0,
    }
    c.SendPacket(basicPkt)
    c.RecvPacket() // receive keepalive response

    // send visithome
    HiLo = utils.Tag2HiLo(tag)
    visitHomeMsg := packets.ClientVisitHome{
        Hi: HiLo[0],
        Lo: HiLo[1],
    }
    basicPkt = packets.Packet{
        Type: packets.MessageType["ClientVisitHome"],
        Version: 0,
        DecryptedPayload: visitHomeMsg.Bytes(),
    }
    c.SendPacket(basicPkt)
    basicPkt = c.RecvPacket() // receive visithome response
    if basicPkt.Type != packets.MessageType["ServerVisitedHome"]{
        log.Printf("ServerVisitedHome not received for %s", tag)
        return basicPkt, errors.New("ServerVisitedHome not received")
    }

    visitHomeResponse := packets.NewServerVisitHomeFromBytes(basicPkt.DecryptedPayload)
    log.Printf("Received VisitHomeData for player %s", utils.HiLo2Tag(visitHomeResponse.Hi, visitHomeResponse.Lo))
    return basicPkt, nil
}
