package config

import (
    "encoding/json"
    "io/ioutil"
    "log"
)

// MyConfig struct
// This is the struct that the config.json must have
type MyConfig struct {
    // RedisDB info
    RedisDB RedisConfig
    // Telegram info
    TelegramBot TelegramConfig
    // CRBot Info
    CRBot CRConfig
}

// DB Config
type RedisConfig struct {
    Host string
    Port string
    Pass string
}

// Telegram Config
type TelegramConfig struct {
    Token string
}

// CR Config
type CRConfig struct {
    PassToken string
    PlayerTag string
    ServerKey string
}

var instance *MyConfig = nil

func CreateInstance(filename string) *MyConfig {
    var err error
    instance, err = loadConfig(filename)
    if err != nil {
        log.Printf("Error loading config file: %s\nUsing default config.", err)
        // use defaults
        instance = &MyConfig{
            RedisDB: RedisConfig{
                Host: "localhost",
                Port: "6379",
                Pass: "",
            },
            TelegramBot: TelegramConfig{
                Token: "",
            },
            CRBot: CRConfig{
                PassToken: "",
                PlayerTag: "",
                ServerKey: "",
            },
        }
    }

    return instance
}

func GetInstance() *MyConfig {
    return instance
}

func loadConfig(filename string) (*MyConfig, error){
    var s *MyConfig

    bytes, err := ioutil.ReadFile(filename)
    if err != nil {
        return s, err
    }
    // Unmarshal json
    err = json.Unmarshal(bytes, &s)
    return s, err
}