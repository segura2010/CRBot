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
    fmt.Printf("CRTGBot v%s\nCommit: %s\nCompilation date: %s\n", CRBOT_VERSION, Commit, CompilationDate)
    fmt.Println("----------------------------------------")
}

func main(){
    cfgFile := flag.String("c", "./config.json", "Config file")
    version := flag.Bool("v", false, "Show current CRBot version")

    if *version{
        showVersionInfo()
    }

}

