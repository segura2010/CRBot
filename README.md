# CRBot: Telegram Bot for Clash Royale player statistics

This is a Telegram Bot that gets CR stats and sends them to you via Telegram.

The bot is splitted in two components:

- CRBot: Clash Royale bot that connects to the CR servers and gets the stats
- TelegramBot: Telegram bot that connects to Telegram servers and listens for update stats requests

When the TelegramBot receives a request, it put a job in the RedisDB server.

The CRBot (you can run multiple bots) will check every X seconds for new jobs in the RedisDB. If there are new jobs, it will get one and it will request the player stats.

Then, the CRBot will send the stats over the RedisDB's channel.

The TelegramBot will listen for completed jobs on the RedisDB's channel, and will respond to the user that requested the stats.

### Installation 

1. Clone this repository and rename the folder to CRBot if it is not the name.
2. Run the install script `install_dependencies.sh` to install all the dependencies.
3. Compile the TelegramBot and CRBot using `install.sh`. Or run `make install` for each one. (IMPORTANT: you should have included your $GOPATH/bin to your $PATH).

### Usage

Start both, the TelegramBot (crtgbot) and CRBot (crbot) binaries. Pass your configuration file path using `-c` flag. Example: `crbot -c /path/to/config.json`.

Check `config_example.json` to create your own config file.

### Why?

I developed the bot just to learn more about Go and because I wanted to create my own tool to check my CR stats.

**Tested on Go 1.8.3**