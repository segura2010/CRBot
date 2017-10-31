package main

import (
    "fmt"
    "flag"
    "time"
    "log"
    "strings"

    "CRBot/db"
    "CRBot/config"
    "CRBot/crbot/bot"
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

    flag.Parse()

    if *version{
        showVersionInfo()
        return
    }

    config.CreateInstance(*cfgFile)

    for{
        job := db.RemovePlayerStatsJob()
        jobInfo := strings.Split(job, ":")
        // log.Printf("playerTag: %s", jobInfo[0])
        if jobInfo[0] != ""{
            visitHomePkt, err := bot.GetStatsForPlayer(jobInfo[0])
            if err == nil{
                // log.Printf("Payload: %x", visitHomePkt.DecryptedPayload)
                msg := fmt.Sprintf("%s:%x", job, visitHomePkt.DecryptedPayload)
                db.PublishToPlayerStats(msg)
            }else{
                log.Printf("ERROR: %s", err)
                msg := fmt.Sprintf("%s:err", job)
                db.PublishToPlayerStats(msg)
            }
        }
        time.Sleep(time.Duration(1) * time.Second)
    }

}

