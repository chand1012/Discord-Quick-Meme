# Discord-Quick-Meme

A discord bot that sends Reddit memes and news to a channel.

## To add

If you are a server owner, just click [here](https://discordapp.com/oauth2/authorize?client_id=438381344943374346&scope=bot) and select your channel to add the bot to your channel.

## To use

- Type `!meme` for a meme from a list of Subreddits, currently [r/funny](https://www.reddit.com/r/funny/), [r/dankmemes](https://www.reddit.com/r/dankmemes/), [r/memes](https://www.reddit.com/r/memes/), and [r/dank_meme](https://www.reddit.com/r/dank_meme/).
- Type `!joke` for a joke from either [r/jokes](https://www.reddit.com/r/jokes/) or [r/darkjokes](https://www.reddit.com/r/darkjokes/).
- Type `!news` for a random news article from [r/news](https://www.reddit.com/r/news/), [r/worldnews](https://www.reddit.com/r/worldnews/), [r/FloridaMan](https://www.reddit.com/r/FloridaMan/), or [r/nottheonion](https://www.reddit.com/r/nottheonion/).
- Type `!5050` or `!fiftyfifty` to pull a post from [r/fiftyfifty](https://reddit.com/r/fiftyfifty).
- A few secret commands if you're willing to look at the source code.
- Type `!news`, `!link`, `!joke`, `!text`, or `!meme` followed by a subreddit name (without the r/) to pull a random top post from that subreddit. The `!text` and `!joke` commands post text directly to the chat. The `!link` and `!news` commands get a post that links a website or webpage. Finally, `!meme` get a piece of media, either a photo, a GIF, or a video. Example: `!joke meanjokes`

## Support

If any support is needed, please post an Issue on the Issues page

## Additional Random Info and Facts

- This bot caches the 100 "hottest" posts on all subreddits called upon it to decrease the response time to the minimum.
- This bot was made as a side project of mine, an idea as some of my friends were not really Reddit users but I always talked about the memes and posts. They use Reddit now.

## Dependencies

- DiscordGo found [discordgo](https://github.com/bwmarrin/discordgo).
- Graw found [graw](https://github.com/turnage/graw).

### Disclaimer

If its offline, I am probably working on it.

Icon by [sandiskplayer34](https://www.deviantart.com/sandiskplayer34) on [DeviantArt](https://www.deviantart.com/sandiskplayer34/art/Reddit-App-Icon-537731823) via [Attribution-ShareAlike 3.0 Unported (CC BY-SA 3.0)](https://creativecommons.org/licenses/by-sa/3.0/).

### To Do

- need to add temporary post blacklist on a per channel basis (using a dictionary will get RAM intensive, maybe a JSON or database)
- get more users
