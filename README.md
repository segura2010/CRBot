# CRBot: Telegram Bot for Clash Royale player statistics

This is a Telegram Bot that gets CR stats and sends them to you via Telegram.

The bot is splitted in two components:

- CRBot: Clash Royale bot that connects to the CR servers and gets the stats
- TelegramBot: Telegram bot that connects to Telegram servers and listens for update stats requests

When the TelegramBot receives a request, it put a job in the RedisDB server.
The CRBot (you can run multiple bots) will check every X seconds for new jobs in the RedisDB. If there are new jobs, it will get one and it will request the player stats.
Then, the CRBot will send the stats over the RedisDB's channel.
The TelegramBot will listen for completed jobs on the RedisDB's channel, and will respond to the user that requested the stats.