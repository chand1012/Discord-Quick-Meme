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

// gets all channel names via the database
func getAllChannelNames() {
	fmt.Println("Getting current channel names and NSFW status...")
	starttime := GetMillis()

	db, err := initDB()

	defer db.Close()

	if err != nil {
		return
	}

	rows, err := db.Query("SELECT channelID, nsfw, name FROM channels")

	if err != nil {
		return
	}

	var channel string
	var nsfwInt int
	var nsfw bool
	var name string

	for rows.Next() {
		err = rows.Scan(&channel, &nsfwInt, &name)

		if nsfwInt == 1 {
			nsfw = true
		} else {
			nsfw = false
		}

		ServerMap[channel] = name
		NSFWMap[channel] = nsfw
	}
	endtime := GetMillis()
	t := endtime - starttime
	fmt.Println("Time to get all current channel names and NSFW status: " + strconv.FormatInt(t, 10) + "ms")
}

// GetChannelData gets channel name and NSFW status
func GetChannelData(discord *discordgo.Session, channelID string, guildID string) (string, bool) {
	fmt.Println("Getting channel data....")

	name, nok := ServerMap[channelID]
	nsfw, wok := NSFWMap[channelID]

	if nok && wok {
		fmt.Println("Values cached.")
		return name, nsfw
	}

	starttime := GetMillis()

	fmt.Println("Checking database....")
	nsfw, name, err := GetChannelFromDB(channelID)

	if err == nil {
		fmt.Println("Channel found in database, adding to RAM....")
		ServerMap[channelID] = name
		NSFWMap[channelID] = nsfw
		endtime := GetMillis()
		t := endtime - starttime
		fmt.Println("Time to get channel name: " + strconv.FormatInt(t, 10) + "ms")
		return name, nsfw
	}

	if err != nil && err != sql.ErrNoRows {
		return channelID, false
	}

	channels, err := discord.GuildChannels(guildID)

	if err != nil {
		fmt.Println("Error getting channel data: ", err)
		return channelID, false
	}

	for _, channel := range channels {
		if channel.ID == channelID {
			ServerMap[channelID] = channel.Name
			NSFWMap[channelID] = channel.NSFW
			endtime := GetMillis()
			t := endtime - starttime
			fmt.Println("Time to get channel name: " + strconv.FormatInt(t, 10) + "ms")
			go AddChannelToDB(channel.ID, channel.NSFW, channel.Name)
			return channel.Name, channel.NSFW
		}
	}

	return channelID, false
}

// gets a channel name from the cache, otherwise searches all channels on server that send the message
func getChannelName(discord *discordgo.Session, channelID string, guildID string) string {
	fmt.Println("Getting channel name....")
	if value, ok := ServerMap[channelID]; ok {
		fmt.Println("Value cached.")
		return value
	}
	starttime := GetMillis()
	channels, err := discord.GuildChannels(guildID)
	if err != nil {
		fmt.Println("Error getting channel name: ", err)
		return channelID
	}
	for _, channel := range channels {
		if channel.ID == channelID {
			ServerMap[channelID] = channel.Name
			endtime := GetMillis()
			t := endtime - starttime
			fmt.Println("Time to get channel name: " + strconv.FormatInt(t, 10) + "ms")
			go AddChannelToDB(channel.ID, channel.NSFW, channel.Name)
			return channel.Name
		}
	}

	return channelID
}

func getChannelNSFW(discord *discordgo.Session, channelID string, guildID string) bool {
	fmt.Println("Getting channel NSFW status....")
	if value, ok := NSFWMap[channelID]; ok {
		fmt.Println("Value cached.")
		return value
	}
	starttime := GetMillis()
	channels, err := discord.GuildChannels(guildID)
	if err != nil {
		fmt.Println("Error getting channel NSFW status: ", err)
		return false
	}
	for _, channel := range channels {
		if channel.ID == channelID {
			NSFWMap[channelID] = channel.NSFW
			endtime := GetMillis()
			t := endtime - starttime
			fmt.Println("Time to get channel NSFW Status: " + strconv.FormatInt(t, 10) + "ms")
			return channel.NSFW
		}
	}

	return false

}

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
func embedSendRoutine(discord *discordgo.Session, channel string, sub string, title string, url string, score int32) {
	rand.Seed(time.Now().Unix())
	randColor := rand.Intn(0xffffff)
	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       randColor,
		Description: "From r/" + sub + "\n Score: " + humanize.Comma(int64(score)),
		Image: &discordgo.MessageEmbedImage{
			URL: url,
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Title:     title,
	}
	_, err := discord.ChannelMessageSendEmbed(channel, embed)
	if err != nil {
		if strings.Contains(err.Error(), string(discordgo.ErrCodeUnknownChannel)) {
			fmt.Println("Channel either was deleted or the bot does not have access to it. Removing all entries from database.")
			err = RemoveChannelFromDBAllTables(channel)
			if err != nil {
				fmt.Println(err)
			}
		}
		fmt.Println("Error posting to channel:", err.Error())
	}
}

// sending a text message routine
func successSendRoutine(discord *discordgo.Session, channel string, sub string, textone string, texttwo string, score int32) {
	_, err := discord.ChannelMessageSend(channel, "From r/"+sub+"\n"+textone+"\n"+texttwo+"\nScore: "+humanize.Comma(int64(score)))
	if err != nil {
		fmt.Println("Error posting to channel:", err.Error())
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

// loop routine that gets posts from the cache.
func getPostLoop(subs []string, postType string, channel string, channelNsfw bool, sort string) (string, int32, string, string, bool, bool) {
	var returnPost QuickPost
	var sub string
	var score int32
	var url string
	var title string
	var nsfw bool

	bannedToggle := false
	toggled := false
	bannedSubs, err := GetAllBannedSubs(channel)
	if err != nil {
		fmt.Println("Error getting banned subreddits: ", err)
	}
	for i := 0; i < 10; i++ {
		returnPost, sub = GetPost(subs, 100, sort, postType)
		blacklisted := CheckBlacklist(channel, returnPost)
		banned := stringInSlice(sub, bannedSubs)

		if !banned {
			LastPost[channel] = returnPost
		}

		score = returnPost.Score
		url = returnPost.Content
		title = returnPost.Title
		nsfw = returnPost.Nsfw

		if nsfw {
			fmt.Println("Post is NSFW.")
		}

		// complicated but it makes sense ish
		if (channelNsfw && !blacklisted && !banned) || (channelNsfw && !nsfw && !blacklisted && !banned) {
			toggled = true
			AddToBlacklist(channel, returnPost)
			break
		} else if !channelNsfw && !nsfw && !blacklisted && !banned {
			toggled = true
			break
		} else {
			if !blacklisted && !banned {
				fmt.Println("Channel is not NSFW but post is NSFW, retrying....")
			} else if banned {
				fmt.Println("Channel banned sub " + sub + ", retrying...")
				if i == 9 {
					bannedToggle = true
				}
			}
		}
	}
	return sub, score, url, title, toggled, bannedToggle

}

// gets a media post and sends it
func getMediaPost(discord *discordgo.Session, channel string, channelNsfw bool, subs []string, sort string) {
	var score int32
	var url string
	var title string
	var sub string
	var bannedToggle bool
	var toggled bool

	imageEndings := []string{".jpg", ".png", ".jpeg"}

	sub, score, url, title, toggled, bannedToggle = getPostLoop(subs, "media", channel, channelNsfw, sort)

	if ContainsAnySubstring(url, imageEndings) && toggled {
		embedSendRoutine(discord, channel, sub, title, url, score)
	} else if toggled {
		successSendRoutine(discord, channel, sub, url, title, score)
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
	var url string
	var title string
	var sub string
	var bannedToggle bool
	var toggled bool

	sub, score, url, title, toggled, bannedToggle = getPostLoop(subs, "link", channel, channelNsfw, sort)

	if toggled {
		successSendRoutine(discord, channel, sub, url, title, score)
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
	if member.User.ID == "" {
		return false
	}
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
