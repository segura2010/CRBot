package main

import (
    "fmt"
    "flag"
    "time"
    "log"

    "CRBot/telegrambot/bot"
    "CRBot/config"
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
    fmt.Printf("CRTGBot v%s\nCommit: %s\nCompilation date: %s\n", CRBOT_VERSION, Commit, CompilationDate)
    fmt.Println("----------------------------------------")
}

func main(){
    cfgFile := flag.String("c", "./config.json", "Config file")
    version := flag.Bool("v", false, "Show current CRBot version")

    flag.Parse()

    if *version{
        showVersionInfo()
        return
    }

    config.CreateInstance(*cfgFile)

    instance := bot.CreateInstance(config.GetInstance().TelegramBot.Token)
    if instance == nil{
        panic("Unable to create TelegramBot")
    }

    log.Printf("Listenning for messages...")
    for{
        // listen for messages...
        time.Sleep(time.Duration(1) * time.Second)
    }

}

