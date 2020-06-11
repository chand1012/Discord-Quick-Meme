# Discord-Quick-Meme

[![Discord Bots](https://top.gg/api/widget/status/438381344943374346.svg)](https://top.gg/bot/438381344943374346) [![Discord Bots](https://top.gg/api/widget/lib/438381344943374346.svg)](https://top.gg/bot/438381344943374346) [![Discord Bots](https://top.gg/api/widget/servers/438381344943374346.svg)](https://top.gg/bot/438381344943374346) ![Build and Test Bot](https://github.com/chand1012/Discord-Quick-Meme/workflows/Build%20and%20Test%20Bot/badge.svg) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=chand1012_Discord-Quick-Meme&metric=alert_status)](https://sonarcloud.io/dashboard?id=chand1012_Discord-Quick-Meme) ![Go Version](https://img.shields.io/github/go-mod/go-version/chand1012/Discord-Quick-Meme/master) [![](https://img.shields.io/discord/626209936262823937)](https://discord.gg/YNnp9uy)

A discord bot that sends Reddit memes and news to a channel.

## To add

If you are a server owner, just click [here](https://discordapp.com/oauth2/authorize?client_id=438381344943374346&scope=bot) and select your channel to add the bot to your channel.

## To use

### General

- Type `!meme` for a meme.
- Type `!joke` for a joke.
- Type `!news` for a random news article.
- Type `!5050` or `!fiftyfifty` to pull a post from [r/fiftyfifty](https://reddit.com/r/fiftyfifty).
- List of subreddits for commands can be found [here](https://github.com/chand1012/Discord-Quick-Meme/blob/master/subs.json).
- A few secret commands if you're willing to look at the source code.
- Type `!news`, `!link`, `!joke`, `!text`, or `!meme` followed by a subreddit name or a list of names separated with spaces (without the r/) to pull a random top post from that subreddit or from a random subreddit from that list. The `!text` and `!joke` commands post text directly to the chat. The `!link` and `!news` commands get a post that links a website or webpage. Finally, `!meme` get a piece of media, either a photo, a GIF, or a video. Examples: `!joke meanjokes` or `!meme cringetopia cringe`
- Type `!source` to get the link to the last sent post in a channel. This was changed because people were complaining memes were coming up twice. This is because Discord shows previews of web content, to remedy this I make the link to the post a requested option, rather than removing it completely.
- Typing `!revsearch` will search the web for the most recently posted PNG or JPEG file in the chat. This is done via Google Reverse Image Search. The search will prioritize results from Reddit, but if a Reddit result cannot be found, Google will return similar images and phrases, which are then put into the chat. **Disclaimer**: Google Reverse Image Search is best at finding simple objects, so don't be surprised if it has issues attempting to find your memes.

### Admins

- If you are a server admin, add a role called `Memebot Admin` to your roles. This role allows you to use the `ban`, `unban`, `subscribe`, and `unsubscribe` commands. 
- The `ban` and `unban` commands allow you ban certain subreddits on either just a channel or the entire server. To ban a subreddit, run the command `!quickmeme ban <mode> <subreddit>` , where mode is either `server` or `channel`. Just replace `ban` with `unban` to do the reverse effect.
    - Typing `!quickmeme getbanned <mode>`, where `<mode>` is either `server` or `channel` will give you a list of the banned subreddits on either the whole server or just the channel the command was executed on. This command can be executed by both admins and regular users.
- (**BETA FEATURE**: May have bugs!) The `subscribe` command allows the bot to periodically post memes in the channel of your choice. The command syntax is `!quickmeme subscribe <interval> <subreddit1>,<subreddit2>,...`. Some examples of intervals are `1h` (hourly), `6h`, `12h`, `1d` (daily), and `1w` (weekly).<sup>There is a maximum time interval of one week and a minimum of 15 minutes.</sup> You can add as many subreddits as you want to the custom command, as long as your command is under 2000 characters (Discord's rule, not mine), and random one will be pulled from your list. The subreddits are separated with _commas_, **not** spaces. If you use spaces you will get an error.
    - To unsubscribe, simply type `!quickmeme unsubscribe` and the bot will stop sending messages in that channel until prompted again.
    - The bot checks the queue every ten seconds, so many of the messages will not be exactly on-the-dot every hour, but should be within a minute or two.

## Support

If any support is needed, please post an Issue on the Issues page on Github or join my support server found [here](https://discord.gg/YNnp9uy).

## Additional Random Info and Facts

- This bot caches the 100 "hottest" posts on all subreddits called upon it to decrease the response time to the minimum.
- This bot was made as a side project of mine, an idea as some of my friends were not really Reddit users but I always talked about the memes and posts. They use Reddit now.

## Dependencies

- [Discordgo](https://github.com/bwmarrin/discordgo).
- [Graw](https://github.com/turnage/graw).
- [go-redis](https://github.com/go-redis/redis).
- [MRISA](https://github.com/vivithemage/mrisa)
- ~~[Redis](https://redis.io/)~~ Migrating to SQL Database

### **Disclaimer**

If its offline, I am probably working on it.

Icon by [sandiskplayer34](https://www.deviantart.com/sandiskplayer34) on [DeviantArt](https://www.deviantart.com/sandiskplayer34/art/Reddit-App-Icon-537731823) via [Attribution-ShareAlike 3.0 Unported (CC BY-SA 3.0)](https://creativecommons.org/licenses/by-sa/3.0/).

### To Do