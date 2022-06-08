# ⚠ Admins

This is the page with all the information for server admins and owners.

## Memebot Admin Role

**ATTENTION: READ THIS BEFORE CONTINUING**

Before you can execute **any** admin-only commands, you must add a role called `Memebot Admin` to your server. This role doesn't have to do anything, just have the name. **Only** give the role to users who should be allowed to change the bot's per server settings, i.e. moderators, people you trust, and yourself. For more information on how to add a role to your Discord server, see [here](https://support.discord.com/hc/en-us/articles/206029707-How-do-I-set-up-Permissions-).

### NSFW Info

As a general Discord rule, the bot **cannot** post NSFW content in a Discord channel that is not marked NSFW. If you want to post explicit and NSFW content from Reddit into your Discord, create a new channel and mark it as NSFW. More details [here](https://support.discord.com/hc/en-us/articles/115000084051-NSFW-channels-and-content).

### Subscriptions

One of the most popular features of the bot is adding a timed meme subscription to your server. This will post a meme on a regular basis to your server in the channel of your choosing. Currently, the smallest allowable time is 15 minutes, and the longest allowable time is one week. In the near future, Patreon Patrons will be allowed to override this limitation.

#### An Example

For this functional example, we will be using the popular subreddit `r/dankmemes`. To subscribe to this channel, with a random post every 15 minutes, simple enter the following (without the `r/`):

`!quickmeme subscribe 15m dankmemes`

The bot will now post memes in your channel every 15 minutes.

![](https://i.imgur.com/yNLe8IL.jpg)

#### More Info

Here is the breakdown of the command:

`!quickmeme subscribe <interval> <subreddit>`

`<interval>` can be any integer followed by `m`, `h`, or `d`, which stand for minutes, hours, and days respectively. For example, a subscription set to `1h` would send memes every hour. A subscriptions set to `45m` would send a meme every 45 minutes.

`<subreddit>` can either be a subreddit or a **comma**-separated (**not** spaces) list of subreddits. If for example you wanted to include `r/memes` and `r/dank_meme` into the previous subscription, you could run the following:

`!quickmeme subscribe 15m dankmemes,memes,dank_meme`

When the 15 minutes are up, the bot will randomly select a post from one of the three subreddits.

To unsubscribe from all memes, simply type `!quickmeme unsubscribe`.

### Subreddit Banning

If for example you _really_ hate a subreddit that your friends keep spamming with the bot, you can ban that subreddit from either the server or the channel. The syntax is:

`!quickmeme ban <mode> <subreddit>`

Where `<mode>` is either `channel` for just a channel ban or `server` to ban a subreddit on the entire Discord server, and `<subreddit>` is the subreddit in question. You can have as many bans as you want, for more bans just repeat the command.

![](https://i.ibb.co/wyZJHLq/image.png)

To unban a subreddit, simply replace the word `ban` with `unban` when executing the command.
