// for any misc discord functions
package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	_ "github.com/go-sql-driver/mysql"
)

// updates the bot status
func updateStatus(discord *discordgo.Session) {
	servers := discord.State.Guilds
	serverCount := int64(len(servers))
	err := discord.UpdateStatus(0, "on "+humanize.Comma(serverCount)+" servers")
	if err != nil {
		fmt.Println("Error updating the status: ", err)
	}
}

// gets a buzzword. See buzz.go
func getBuzzWord(discord *discordgo.Session, channel string) error {
	statement := getABuzzWord()
	_, err := discord.ChannelMessageSend(channel, statement)
	return err
}

// send routine for embedded messages
func embedSendRoutine(discord *discordgo.Session, channel string, sub string, title string, contentURL string, score int32) {
	rand.Seed(time.Now().Unix())
	randColor := rand.Intn(0xffffff)
	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       randColor,
		Description: "From r/" + sub + "\n Score: " + humanize.Comma(int64(score)),
		Image: &discordgo.MessageEmbedImage{
			URL: contentURL,
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Title:     title,
	}
	_, err := discord.ChannelMessageSendEmbed(channel, embed)
	if err != nil {
		if strings.Contains(err.Error(), fmt.Sprint(discordgo.ErrCodeUnknownChannel)) {
			fmt.Println("Channel either was deleted or the bot does not have access to it. Removing all entries from database.")
			err = RemoveChannelFromDBAllTables(channel)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Error posting to channel:", err.Error())
		}
	}
}

// sending a text message routine
func successSendRoutine(discord *discordgo.Session, channel string, sub string, textone string, texttwo string, score int32) {
	_, err := discord.ChannelMessageSend(channel, "From r/"+sub+"\n"+textone+"\n"+texttwo+"\nScore: "+humanize.Comma(int64(score)))
	if err != nil {
		if strings.Contains(err.Error(), fmt.Sprint(discordgo.ErrCodeUnknownChannel)) {
			fmt.Println("Channel either was deleted or the bot does not have access to it. Removing all entries from database.")
			err = RemoveChannelFromDBAllTables(channel)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Error posting to channel:", err.Error())
		}
	}
}

// banned subreddit routine
func bannedSendRoutine(discord *discordgo.Session, channel string, sub string) {
	discord.ChannelMessageSend(channel, "Error!\nToo many attempts due to find an acceptable post due to a banned subreddit: "+sub)
}

// nsfw subreddit not allowed routine
func nsfwSendRoutine(discord *discordgo.Session, channel string) {
	discord.ChannelMessageSend(channel, "Error!\nToo many tries to not find NSFW post, maybe that Subreddit is filled with them? Hint: Name sure that the channel is marked as \"NSFW\".")
}

func errSendRoutine(discord *discordgo.Session, channel string, err error) {
	discord.ChannelMessageSend(channel, "Critical Error!\nThere was a critical error. "+err.Error()+" Please report this if possible to the Github page: https://github.com/chand1012/Discord-Quick-Meme/issues")
}

// gets a media post and sends it
func getMediaPost(discord *discordgo.Session, channel string, channelNsfw bool, subs []string, sort string) {
	var score int32
	var contentURL string
	var title string
	var sub string
	var bannedToggle bool
	var toggled bool

	imageEndings := []string{".jpg", ".png", ".jpeg"}

	sub, score, contentURL, title, toggled, bannedToggle = getPostLoop(subs, "media", channel, channelNsfw, sort)

	if ContainsAnySubstring(contentURL, imageEndings) && toggled {
		embedSendRoutine(discord, channel, sub, title, contentURL, score)
	} else if toggled {
		successSendRoutine(discord, channel, sub, contentURL, title, score)
	} else if bannedToggle {
		bannedSendRoutine(discord, channel, sub)
	} else {
		nsfwSendRoutine(discord, channel)
	}

}

func getMediaPostSettings(discord *discordgo.Session, channel string, channelNsfw bool, subs []string, sort string, settings guildSettings) {
	var score int32
	var contentURL string
	var title string
	var sub string
	var bannedToggle bool
	var toggled bool

	imageEndings := []string{".jpg", ".png", ".jpeg"}
	videoEndings := []string{".mp4"}

	sub, score, contentURL, title, toggled, bannedToggle = getPostLoop(subs, "media", channel, channelNsfw, sort)

	if ContainsAnySubstring(contentURL, imageEndings) && toggled {
		if settings.Proxy {
			if settings.ProxyMode == 1 {
				uploadEmbedSendRoutine(discord, channel, sub, title, contentURL, score)
			} else {
				proxyEmbedSendRoutine(discord, channel, sub, title, contentURL, score)
			}
		} else {
			embedSendRoutine(discord, channel, sub, title, contentURL, score)
		}
	} else if ContainsAnySubstring(contentURL, videoEndings) && toggled {
		if settings.Proxy {
			if settings.ProxyMode == 1 {
				uploadSendRoutine(discord, channel, sub, title, contentURL, score)
			} else {
				proxySendRoutine(discord, channel, sub, title, contentURL, score)
			}
		} else {
			successSendRoutine(discord, channel, sub, contentURL, title, score)
		}
	} else if toggled {
		successSendRoutine(discord, channel, sub, contentURL, title, score)
	} else if bannedToggle {
		bannedSendRoutine(discord, channel, sub)
	} else {
		nsfwSendRoutine(discord, channel)
	}
}

// gets a text post and sends it
func getTextPost(discord *discordgo.Session, channel string, channelNsfw bool, subs []string, sort string) {
	var score int32
	var text string
	var title string
	var sub string
	var bannedToggle bool
	var toggled bool

	sub, score, text, title, toggled, bannedToggle = getPostLoop(subs, "text", channel, channelNsfw, sort)

	if toggled {
		successSendRoutine(discord, channel, sub, title, text, score)
	} else if bannedToggle {
		bannedSendRoutine(discord, channel, sub)
	} else {
		nsfwSendRoutine(discord, channel)
	}
}

// gets a link post and sends it
func getLinkPost(discord *discordgo.Session, channel string, channelNsfw bool, subs []string, sort string) {
	var score int32
	var contentURL string
	var title string
	var sub string
	var bannedToggle bool
	var toggled bool

	sub, score, contentURL, title, toggled, bannedToggle = getPostLoop(subs, "link", channel, channelNsfw, sort)

	if toggled {
		successSendRoutine(discord, channel, sub, contentURL, title, score)
	} else if bannedToggle {
		bannedSendRoutine(discord, channel, sub)
	} else {
		nsfwSendRoutine(discord, channel)
	}
}

// gets the number of users across all servers with the bot
// Discordgo broke this, gonna keep it here in case it gets fixed
// func getNumberOfUsers(discord *discordgo.Session) int {
// 	count := 0
// 	for _, guild := range discord.State.Guilds {
// 		count += len(guild.Members)
// 	}
// 	return count
// }

// gets the user's member struct via their
// func getUserMemberFromGuild(discord *discordgo.Session, guildID string, user discordgo.User) discordgo.Member {
// 	guildObject, _ := discord.Guild(guildID)
// 	for _, member := range guildObject.Members {
// 		if member.User.ID == user.ID {
// 			return *member
// 		}
// 	}
// 	return discordgo.Member{}
// }

// checks if user is a memebot admin.
func isUserMemeBotAdmin(discord *discordgo.Session, guildID string, user *discordgo.User) bool {
	adminCode := "memebot admin"
	member, _ := discord.GuildMember(guildID, user.ID)
	guildRoles, _ := discord.GuildRoles(guildID)
	for _, role := range guildRoles {
		for _, roleID := range member.Roles {
			if role.ID == roleID && strings.Contains(strings.ToLower(role.Name), adminCode) {
				return true
			}
		}
	}
	return false
}

// gets the source of the last sent meme
func getSource(discord *discordgo.Session, channel string) error {
	var err error
	returnPost := LastPost[channel]
	if returnPost.Permalink == "" {
		_, err = discord.ChannelMessageSend(channel, "ERROR: Nothing has ever been sent on this channel!")
	} else {
		postlink := returnPost.Permalink
		_, err = discord.ChannelMessageSend(channel, "Source: https://reddit.com"+postlink)
	}
	return err
}

// ban a subreddit routine
func banSubRoutine(discord *discordgo.Session, channel string, commandContent []string, guildID string, user *discordgo.User) {
	if len(commandContent) != 4 {
		discord.ChannelMessageSend(channel, "Incorrect command syntax! Correct syntax is `!quickmeme ban [mode] [subreddit]`\nMode can be `channel` or `server`.")
	} else if isUserMemeBotAdmin(discord, guildID, user) { // fix this
		if commandContent[2] == "server" {
			channels, _ := discord.GuildChannels(guildID)
			subreddits := textFilterSlice(commandContent[3:])
			if subreddits == nil {
				discord.ChannelMessageSend(channel, ErrorMsg)
				return
			}
			for _, chat := range channels {
				// this should be async to save time
				for _, subreddit := range subreddits {
					go SetBannedSubreddit(chat.ID, strings.ToLower(subreddit))
				}
			}
			// this should be a message about the ban
			discord.ChannelMessageSend(channel, user.Mention()+" banned subreddit(s) "+strings.Join(subreddits, ", ")+" on all channels.")
		} else {
			subreddits := textFilterSlice(commandContent[3:])
			if subreddits == nil {
				discord.ChannelMessageSend(channel, ErrorMsg)
				return
			}
			for _, subreddit := range subreddits {
				go SetBannedSubreddit(channel, subreddit)
			}
			discord.ChannelMessageSend(channel, user.Mention()+" banned subreddit(s) "+strings.Join(subreddits, ", ")+".")
		}
	} else {
		discord.ChannelMessageSend(channel, ErrorMsg)
	}
}

// unban a subreddit routine
func unbanSubRoutine(discord *discordgo.Session, channel string, commandContent []string, guildID string, user *discordgo.User) {
	if len(commandContent) < 4 || len(commandContent) > 5 {
		discord.ChannelMessageSend(channel, "Incorrect command syntax! Correct syntax is `!quickmeme unban [mode] [subreddit]`\nMode can be `channel` or `server`.")
	} else if isUserMemeBotAdmin(discord, guildID, user) { // fix this
		if commandContent[2] == "server" {
			channels, _ := discord.GuildChannels(guildID)
			subreddits := textFilterSlice(commandContent[3:])
			if subreddits == nil {
				discord.ChannelMessageSend(channel, ErrorMsg)
				return
			}
			for _, chat := range channels {
				// this should be async to save time
				for _, subreddit := range subreddits {
					go RemoveBannedSubreddit(chat.ID, subreddit)
				}
			}
			// this should be a message about the ban
			discord.ChannelMessageSend(channel, user.Mention()+" unbanned subreddit(s) "+strings.Join(subreddits, ", ")+" on all channels.")
		} else {
			// there should be a message about the ban here
			subreddits := textFilterSlice(commandContent[3:])
			if subreddits == nil {
				discord.ChannelMessageSend(channel, ErrorMsg)
				return
			}
			for _, subreddit := range subreddits {
				go RemoveBannedSubreddit(channel, subreddit)
			}
			discord.ChannelMessageSend(channel, user.Mention()+" unbanned subreddit(s) "+strings.Join(subreddits, ", ")+".")
		}
	} else {
		discord.ChannelMessageSend(channel, "Insufficient Permissions! You must have the \"Memebot Admin\" role to ban subreddits!")
	}
}

// gets banned subreddits and send to channel
func getbannedSubRoutine(discord *discordgo.Session, channel string, commandContent []string, guildID string, user *discordgo.User) {
	var banContext string
	if len(commandContent) < 3 {
		banContext = "channel"
	} else {
		banContext = commandContent[2]
	}

	if banContext == "server" {
		channels, _ := discord.GuildChannels(guildID)
		for _, chat := range channels {
			bannedSubs, err := GetAllBannedSubs(chat.ID)
			if err != nil {
				discord.ChannelMessageSend(channel, ErrorMsg)
				fmt.Println(err)
				break
			}
			msgString := strings.Join(bannedSubs, ", ")
			if msgString != "" && chat.Type == discordgo.ChannelTypeGuildText {
				discord.ChannelMessageSend(channel, "Subs banned on "+chat.Name+":\n"+msgString)
			}
		}
	} else {
		bannedSubs, err := GetAllBannedSubs(channel)
		if err != nil {
			discord.ChannelMessageSend(channel, ErrorMsg)
			fmt.Println(err)
		} else {
			msgString := strings.Join(bannedSubs, ", ")
			discord.ChannelMessageSend(channel, "Subs banned on this channel:\n"+msgString)
		}
	}

}

func setQueueRoutine(discord *discordgo.Session, channel string, commandContent []string, channelNsfw bool) {
	var err error
	if len(commandContent) == 4 {
		interval := commandContent[2]
		expressions := []string{"[0-9]+m", "[0-9]+h", "[0-9]+d"}
		match := matchRegexList(expressions, interval)

		if !match {
			discord.ChannelMessageSend(channel, "There was an error with your syntax, see here: https://bit.ly/DiscordQuickMemeAdminSyntax")
			return
		}

		err = SetMemeQueue(channel, channelNsfw, interval, commandContent[3])
		if err != nil {
			fmt.Println("There was an error setting the Queue: " + err.Error())
			errSendRoutine(discord, channel, err)
			return
		}
		discord.ChannelMessageSend(channel, "Memes will now be sent at a regularly scheduled time.")
	} else {
		discord.ChannelMessageSend(channel, "There was an error with your syntax, see here: https://bit.ly/DiscordQuickMemeAdminSyntax")
		return
	}
}

func delQueueRoutine(discord *discordgo.Session, channel string) {
	fmt.Println("Deleting channel", channel)
	err := DeleteMemeQueue(channel)
	fmt.Println("Checking for errors...")
	if err == sql.ErrNoRows {
		discord.ChannelMessageSend(channel, "Error, this channel isn't subscribed.")
		return
	}
	if err != nil {
		fmt.Println("There was an error removing an item from the queue: ", err)
		errSendRoutine(discord, channel, err)
		return
	}
	discord.ChannelMessageSend(channel, "You have successfully unsubscribed this channel from memes.")
	fmt.Println("Done.")
}

func helpRoutine(discord *discordgo.Session, channel string) {
	discord.ChannelMessageSend(channel, "For command syntax, see here: https://github.com/chand1012/Discord-Quick-Meme#to-use")
}

func updateProxyRoutine(discord *discordgo.Session, channel string, guildID string, commandContent []string, settings guildSettings) {
	var setting string
	var value string
	var proxyEnable bool
	var proxyMode int8

	if len(commandContent) != 4 {
		discord.ChannelMessageSend(channel, "Incorrect command syntax! Correct syntax is `!quickmeme proxy <setting> <value>` see here for more info: https://github.com/chand1012/Discord-Quick-Meme#to-use")
		return
	}

	proxyEnable = settings.Proxy
	proxyMode = settings.ProxyMode
	setting = commandContent[2]
	value = commandContent[3]

	switch setting {
	case "enable":
		if value == "false" {
			proxyEnable = false
			discord.ChannelMessageSend(channel, "Disabling Proxy.")
		} else {
			proxyEnable = true
			discord.ChannelMessageSend(channel, "Enabling Proxy.")
		}
	case "mode":
		if value == "discord" {
			proxyMode = 0
			discord.ChannelMessageSend(channel, "Setting proxy mode to Discord Proxy.")
		} else if value == "worker" {
			proxyMode = 1
			discord.ChannelMessageSend(channel, "Setting proxy mode to Custom Proxy.")
		} else {
			discord.ChannelMessageSend(channel, "Incorrect command syntax! Proxy Mode must be `discord` for Discord's image upload as a proxy or `worker` for the external proxy.")
			return
		}
	}

	err := SetGuildStatus(guildID, proxyEnable, proxyMode)

	if err != nil {
		fmt.Println(err)
		errSendRoutine(discord, channel, err)
	} else {
		discord.ChannelMessageSend(channel, "Proxy settings updated.")
	}

}

func setBenefitsRoutine(discord *discordgo.Session, channel string, guildID string, user *discordgo.User) {
	userStatus, err := getPatronStatus(user.ID)

	if err != nil {
		fmt.Println(err)
		errSendRoutine(discord, channel, err)
		return
	}

	if userStatus == 0 {
		discord.ChannelMessageSend(channel, user.Mention()+" , you are not a Patron! Subscribe and get some awesome benefits here: https://www.patreon.com/DiscordQuickMeme")
		return
	}

	userStatus, cooldown, err := getBenefitServer(user.ID, guildID)

	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		errSendRoutine(discord, channel, err)
		return
	}

	if cooldown != 0 || err != sql.ErrNoRows {
		discord.ChannelMessageSend(channel, "This server is already enrolled in our Patreon benefits, silly!")
		return
	}

	userStatus, benefitGuilds, err := getAllBenefitsForUser(user.ID)

	if userStatus == 1 && len(benefitGuilds) > 0 {
		discord.ChannelMessageSend(channel, user.Mention()+" has already met their benefit server limit. Either upgrade your package to allow for more servers, remove the server you have already benefit, or have someone else give this server benefits: https://www.patreon.com/DiscordQuickMeme")
		return
	} else if userStatus == 2 && len(benefitGuilds) >= 3 {
		discord.ChannelMessageSend(channel, user.Mention()+" has already met their benefit server limit. Either remove a server from your benefits, or have someone else give this server benefits: https://www.patreon.com/DiscordQuickMeme")
		return
	}

	err = setBenefitServer(user.ID, userStatus, guildID)

	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		errSendRoutine(discord, channel, err)
		return
	}

	discord.ChannelMessageSend(channel, "Hey, @everyone ! "+user.Mention()+" just gave you QuickMeme server benefits! Give them love! :clap:")
}

func removeBenefitsRoutine(discord *discordgo.Session, channel string, guildID string, user *discordgo.User) {

	_, cooldown, err := getBenefitServer(user.ID, guildID)

	if err == sql.ErrNoRows {
		discord.ChannelMessageSend(channel, "This server isn't subscribed to QuickMeme benefits. If you want to subscribe, please sign up here: https://www.patreon.com/DiscordQuickMeme")
		return
	} else if err != nil {
		fmt.Println(err)
		errSendRoutine(discord, channel, err)
		return
	}

	if cooldown < time.Now().Unix() {
		waitTime := (cooldown - time.Now().Unix()) / 86400 // seconds in a day
		discord.ChannelMessageSend(channel, "You must wait 30 days before changing your server. You have about "+strconv.FormatInt(waitTime, 10)+" days before you can change servers.")
		return
	}

	err = removeBenefitServer(guildID)

	if err != nil {
		fmt.Println(err)
		errSendRoutine(discord, channel, err)
		return
	}

	discord.ChannelMessageSend(channel, user.Mention()+" has removed their benefits from this server. You may now take your benefits elsewhere. If anyone else would like to provide QuickMeme benefits to this server, sign up can be found here: https://www.patreon.com/DiscordQuickMeme")
}
