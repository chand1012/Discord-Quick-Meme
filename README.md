# Discord-Quick-Meme
A discord bot that sends Reddit memes and news to a channel.
## To add
If you are a server owner, just click [here](https://discordapp.com/oauth2/authorize?client_id=438381344943374346&scope=bot) and select your channel to add the bot to your channel.

## To use
- Type `!meme` for a meme from a list of subreddits, currently [r/funny](https://www.reddit.com/r/funny/), [r/dankmemes](https://www.reddit.com/r/dankmemes/), [r/memes](https://www.reddit.com/r/memes/), and [r/dank_meme](https://www.reddit.com/r/dank_meme/).
- Type `!joke` for a joke from either [r/jokes](https://www.reddit.com/r/jokes/) or [r/darkjokes](https://www.reddit.com/r/darkjokes/).
- Type `news` for a random news article from [r/news](https://www.reddit.com/r/news/), [r/worldnews](https://www.reddit.com/r/worldnews/), [r/FloridaMan](https://www.reddit.com/r/FloridaMan/), or [r/nottheonion](https://www.reddit.com/r/nottheonion/).
- Type `!news`, `!joke`, or `!meme` followed by a subreddit name (without the r/) to pull a random top post from that subreddit. Example: `!joke meanjokes`
- If the meme or joke is nsfw, it will not be used on a regular channel. The channel must have `nsfw` in the title to allow for the post to go through. Otherwise, it spits out an error.

## Dependencies
- discord.py found [here](https://github.com/Rapptz/discord.py/).
- Praw found [here](https://github.com/praw-dev/praw).
- the included lib file

### Disclaimer
I AM NOT RESPONSIBLE FOR ANY CONTENT THAT IS POSTED WITH THIS BOT. IT IS AT THE DISCRETION OF THE SERVER OWNER AND USERS.

This bot is under heavy development and maintenance, so it may not be online all the time. If its offline, I am probably working on it.

### To do
- make the code simpler and therefore faster
- migrate the bot from a Raspberry Pi to an actual server
- add more features
