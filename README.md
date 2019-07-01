
# Discord-Quick-Meme
A discord bot that sends Reddit memes and news to a channel.

## To add
If you are a server owner, just click [here](https://discordapp.com/oauth2/authorize?client_id=438381344943374346&scope=bot) and select your channel to add the bot to your channel.

## To use
- Type `!meme` for a meme from a list of Subreddits, currently [r/funny](https://www.reddit.com/r/funny/), [r/dankmemes](https://www.reddit.com/r/dankmemes/), [r/memes](https://www.reddit.com/r/memes/), and [r/dank_meme](https://www.reddit.com/r/dank_meme/).
- Type `!joke` for a joke from either [r/jokes](https://www.reddit.com/r/jokes/) or [r/darkjokes](https://www.reddit.com/r/darkjokes/).
- Type `!news` for a random news article from [r/news](https://www.reddit.com/r/news/), [r/worldnews](https://www.reddit.com/r/worldnews/), [r/FloridaMan](https://www.reddit.com/r/FloridaMan/), or [r/nottheonion](https://www.reddit.com/r/nottheonion/).
- Type `!5050` or `!fiftyfifty` to pull a post from [r/fiftyfifty](https://reddit.com/r/fiftyfifty).
- Type `!post` followed by a post ID to directly pull a post from reddit. Example: `!post 8f3wan`.
- A few secret commands if you're willing to look at the source code.
- Type `!news`, `!joke`, or `!meme` followed by a subreddit name (without the r/) to pull a random top post from that subreddit. Example: `!joke meanjokes`
- If the meme or joke is nsfw, it will not be used on a regular channel. The channel must have `nsfw` in the title to allow for the post to go through. If the bot cannot find a SFW post on the subreddit given, it spits out an error.

## Support
If any support is needed, please post an Issue on the Issues page. 

## Additional Random Info and Facts
- The bot will not allow you to post more than 5 memes a minute per channel. This was a preventative measure so that the ~~VPS~~ server the bot is running on does not get overloaded.
- This bot was made as a side project of mine, an idea as some of my friends were not really Reddit users but I always talked about the memes and posts. They use Reddit now.

## Dependencies
- discord.py rewrite found [here](https://github.com/Rapptz/discord.py/).
- Praw found [here](https://github.com/praw-dev/praw).

### Disclaimer

If its offline, I am probably working on it.

Icon by [sandiskplayer34](https://www.deviantart.com/sandiskplayer34) on [DeviantArt](https://www.deviantart.com/sandiskplayer34/art/Reddit-App-Icon-537731823). 

### To Do
- need to add temporary post blacklist on a per channel basis (using a dictionary will get RAM intensive, maybe a JSON or database)
- get more users
