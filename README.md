# Discord-Quick-Meme

[![Discord Bots](https://top.gg/api/widget/status/438381344943374346.svg)](https://top.gg/bot/438381344943374346) [![Discord Bots](https://top.gg/api/widget/lib/438381344943374346.svg)](https://top.gg/bot/438381344943374346) ![](https://github.com/chand1012/Discord-Quick-Meme/workflows/Go/badge.svg) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=chand1012_Discord-Quick-Meme&metric=alert_status)](https://sonarcloud.io/dashboard?id=chand1012_Discord-Quick-Meme)

<a href="https://www.buymeacoffee.com/chand1012" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/lato-blue.png" alt="Buy Me A Coffee" style="height: 51px !important;width: 217px !important;" ></a>

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
- Typing `!search` will search the web for the most recently posted PNG or JPEG file in the chat. This is done via Google Reverse Image Search. The search will prioritize results from Reddit, but if a Reddit result cannot be found, Google will return similar images and phrases, which are then put into the chat. **Disclaimer**: Google Reverse Image Search is best at finding simple objects, so don't be surprised if it has issues attempting to find your memes.

### Admins

- If you are a server admin, add a role called `Memebot Admin` to your roles. This role allows you to use the `ban` and `unban` commands. These commands allow you ban certain subreddits on either just a channel or the entire server. To ban a subreddit, run the command `!quickmeme ban <mode> <subreddit>` to ban a subreddit, where mode is either `server` or `channel`.
- Typing `!quickmeme getbanned <mode>`, where `<mode>` is either `server` or `channel` will give you a list of the banned subreddits on either the whole server or just the channel the command was executed on. This command can be executed by both admins and regular users.

## Support

If any support is needed, please post an Issue on the Issues page on Github or join my support server found [here](https://discord.gg/YNnp9uy).

## Additional Random Info and Facts

- This bot caches the 100 "hottest" posts on all subreddits called upon it to decrease the response time to the minimum.
- This bot was made as a side project of mine, an idea as some of my friends were not really Reddit users but I always talked about the memes and posts. They use Reddit now.

## Dependencies

- [Discordgo](https://github.com/bwmarrin/discordgo).
- [Graw](https://github.com/turnage/graw).
- [go-redis](https://github.com/go-redis/redis).
- [MRISA](https://github.com/vivithemage/mrisa) running on the same local machine as the bot, or with the address correctly updated in `data.json`.
- [Redis](https://redis.io/) running on the same local machine as the bot, or with the address correctly updated in `data.json`.

### **Disclaimer**

If its offline, I am probably working on it.

Icon by [sandiskplayer34](https://www.deviantart.com/sandiskplayer34) on [DeviantArt](https://www.deviantart.com/sandiskplayer34/art/Reddit-App-Icon-537731823) via [Attribution-ShareAlike 3.0 Unported (CC BY-SA 3.0)](https://creativecommons.org/licenses/by-sa/3.0/).

### Coming Soon

- Refactor and Further optimization for faster memes.

### To Do
- Refactor and organize better (especially the main loop)
- get more users

