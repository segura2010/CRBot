package main

import (
    "fmt"
    "encoding/hex"
    "flag"

    "github.com/segura2010/cr-go/client"
    "github.com/segura2010/cr-go/packets"
    "github.com/segura2010/cr-go/utils"
)

/* Update version number on each release:
    Given a version number x.y.z, increment the:

    x - major release
    y - minor release
    z - build number
*/
const CRBOT_VERSION = "0.0.1"
var Commit string
var CompilationDate string

func showVersionInfo(){
    fmt.Println("----------------------------------------")
    fmt.Printf("CRBot v%s\nCommit: %s\nCompilation date: %s\n", CRBOT_VERSION, Commit, CompilationDate)
    fmt.Println("----------------------------------------")
}

func main(){
    cfgFile := flag.String("c", "./config.json", "Config file")
    version := flag.Bool("v", false, "Show current CRBot version")

    if *version{
        showVersionInfo()
    }

    serverKey_202 := "980cf7bb7262b386fea61034aba7370613627919666b34e6ecf66307a381dd61"
    serverKey,_ := hex.DecodeString(serverKey_202)

    helloPayload := packets.NewDefaultClientHello()
    helloPkt := packets.Packet{
        Type: packets.MessageType["ClientHello"],
        Version: 0,
        Payload: helloPayload.Bytes(),
    }

    loginPayload := packets.NewDefaultClientLogin()
    HiLo := utils.Tag2HiLo("") // test account
    loginPayload.Hi = HiLo[0]
    loginPayload.Lo = HiLo[1]
    loginPayload.PassToken = ""

    loginPkt := packets.Packet{
        Type: packets.MessageType["ClientLogin"],
        Version: 0,
        DecryptedPayload: loginPayload.Bytes(),
    }

    var basicPkt packets.Packet

    serverAddress := "0.0.0.0:9339"
    //serverAddress := "game.clashroyaleapp.com:9339"

    c := client.NewCRClient(serverKey)
    c.Connect(serverAddress)
    defer c.Close()

    c.SendPacket(helloPkt)
    c.RecvPacket() // receive hello response

    c.SendPacket(loginPkt)
    basicPkt = c.RecvPacket() // receive login response
    loginOk := packets.NewServerLoginOkFromBytes(basicPkt.DecryptedPayload)
    fmt.Printf("\nLoginOk: %v", loginOk)

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
    HiLo = utils.Tag2HiLo("UY800JJ")
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
    //fmt.Printf("\nVisitedHomeRequest: %x", visitHomeMsg.Bytes())
    //fmt.Printf("\nVisitedHome: %x", basicPkt.DecryptedPayload)

    visitHomeResponse := packets.NewServerVisitHomeFromBytes(basicPkt.DecryptedPayload)
    fmt.Printf("\nVisitHomeMsg: %x", visitHomeMsg.Bytes())
    fmt.Printf("\n%v", visitHomeResponse)
    fmt.Printf("\n%s", utils.HiLo2Tag(visitHomeResponse.Hi, visitHomeResponse.Lo))

}

