// for any misc discord functions
package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Server server object for the golang channels
type Server struct {
	IDs   []string
	Names []string
}

func getAllWorker(discord *discordgo.Session, guildID string, send chan<- Server, wg *sync.WaitGroup, workerNumber int) {
	defer wg.Done()
	var ids []string
	var names []string
	channels, _ := discord.GuildChannels(guildID)
	for _, channel := range channels {
		if channel.Type != discordgo.ChannelTypeGuildText {
			continue
		}
		ids = append(ids, channel.ID)
		names = append(names, channel.Name)
	}
	server := Server{
		IDs:   ids,
		Names: names,
	}
	send <- server
}

func getAllChannelNames(discord *discordgo.Session) {
	var wg sync.WaitGroup
	fmt.Println("Getting current channel names...")
	starttime := GetMillis()
	guilds := discord.State.Guilds
	bufferSize := len(guilds)
	recv := make(chan Server, bufferSize)
	for i, guild := range guilds {
		wg.Add(1)
		go getAllWorker(discord, guild.ID, recv, &wg, i)
	}
	wg.Wait()
	close(recv)
	for i := 0; i < bufferSize; i++ {
		thing := <-recv
		length := len(thing.IDs)
		for x := 0; x < length; x++ {
			ServerMap[thing.IDs[x]] = thing.Names[x]
		}
	}
	endtime := GetMillis()
	t := endtime - starttime
	fmt.Println("Time to get all current channel names: " + strconv.FormatInt(t, 10) + "ms")
}

func getChannelName(discord *discordgo.Session, channelid string, guildID string) string {
	fmt.Println("Getting channel name....")
	if _, ok := ServerMap[channelid]; ok {
		return ServerMap[channelid]
	}
	starttime := GetMillis()
	channels, _ := discord.GuildChannels(guildID)
	for _, channel := range channels {
		if channel.ID == channelid {
			ServerMap[channelid] = channel.Name
			endtime := GetMillis()
			t := endtime - starttime
			fmt.Println("Time to get channel name: " + strconv.FormatInt(t, 10) + "ms")
			return channel.Name
		}
	}

	return ""
}

func updateStatus(discord *discordgo.Session) {
	uCount := getNumberOfUsers(discord)
	err := discord.UpdateStatus(0, "with "+strconv.Itoa(uCount)+" others")
	if err != nil {
		errCheck("Error attempting to set the status.", err)
	}
}

// ComesFromDM returns true if a message comes from a DM channel
func ComesFromDM(s *discordgo.Session, m *discordgo.MessageCreate) (bool, error) {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		if channel, err = s.Channel(m.ChannelID); err != nil {
			return false, err
		}
	}

	return channel.Type == discordgo.ChannelTypeDM, nil
}

func getBuzzWord(discord *discordgo.Session, channel string) error {
	statement := getABuzzWord()
	_, err := discord.ChannelMessageSend(channel, statement)
	return err
}

func successSendRoutine(discord *discordgo.Session, channel string, sub string, textone string, texttwo string, score int32) {
	discord.ChannelMessageSend(channel, "From r/"+sub)
	discord.ChannelMessageSend(channel, textone)
	discord.ChannelMessageSend(channel, texttwo)
	discord.ChannelMessageSend(channel, "Score: "+strconv.FormatInt(int64(score), 10))
}

func bannedSendRoutine(discord *discordgo.Session, channel string) {
	discord.ChannelMessageSend(channel, "Error!")
	discord.ChannelMessageSend(channel, "Too many tries to find a post on an unbanned subreddit!")
}

func nsfwSendRoutine(discord *discordgo.Session, channel string) {
	discord.ChannelMessageSend(channel, "Error!")
	discord.ChannelMessageSend(channel, "Too many tries to not find NSFW post, maybe that Subreddit is filled with them? Hint: Name sure that the channel is marked as \"NSFW\".")
}

func getPostLoop(subs []string, postType string, channel string, channelNsfw bool, sort string) (string, int32, string, string, bool, bool) {
	var returnPost QuickPost
	var sub string
	var score int32
	var url string
	var title string
	var nsfw bool
	bannedToggle := false
	toggled := false
	bannedSubs, _ := GetBannedSubreddits(channel)
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
		if channelNsfw && !blacklisted && !banned {
			toggled = true
			AddToBlacklist(channel, returnPost)
			break
		} else if channelNsfw && !nsfw && !blacklisted && !banned {
			toggled = true
			AddToBlacklist(channel, returnPost)
			break
		} else if !channelNsfw && !nsfw && !blacklisted && !banned {
			toggled = true
			break
		} else {
			if !blacklisted && !banned {
				fmt.Println("Channel is not NSFW but post is NSFW, retrying...")
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

func getMediaPost(discord *discordgo.Session, channel string, channelNsfw bool, subs []string, sort string) {
	var score int32
	var url string
	var title string
	var sub string
	var bannedToggle bool
	var toggled bool

	rand.Seed(time.Now().Unix())
	randColor := rand.Intn(0xffffff)
	imageEndings := []string{".jpg", ".png", ".jpeg"}

	sub, score, url, title, toggled, bannedToggle = getPostLoop(subs, "media", channel, channelNsfw, sort)

	if ContainsAnySubstring(url, imageEndings) && toggled {
		embed := &discordgo.MessageEmbed{
			Author:      &discordgo.MessageEmbedAuthor{},
			Color:       randColor,
			Description: "Score: " + strconv.FormatInt(int64(score), 10),
			Image: &discordgo.MessageEmbedImage{
				URL: url,
			},
			Timestamp: time.Now().Format(time.RFC3339),
			Title:     title,
		}
		discord.ChannelMessageSend(channel, "From r/"+sub)
		discord.ChannelMessageSendEmbed(channel, embed)
	} else if toggled {
		successSendRoutine(discord, channel, sub, url, title, score)
	} else if bannedToggle {
		bannedSendRoutine(discord, channel)
	} else {
		nsfwSendRoutine(discord, channel)
	}

}

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
		bannedSendRoutine(discord, channel)
	} else {
		nsfwSendRoutine(discord, channel)
	}
}

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
		bannedSendRoutine(discord, channel)
	} else {
		nsfwSendRoutine(discord, channel)
	}
}

func getNumberOfUsers(discord *discordgo.Session) int {
	count := 0
	for _, guild := range discord.State.Guilds {
		count += len(guild.Members)
	}
	return count
}

// gets the user's member struct via their
func getUserMemberFromGuild(discord *discordgo.Session, guildID string, user discordgo.User) discordgo.Member {
	guildObject, _ := discord.Guild(guildID)
	for _, member := range guildObject.Members {
		if member.User.ID == user.ID {
			return *member
		}
	}
	return discordgo.Member{}
}

func isUserMemeBotAdmin(discord *discordgo.Session, guildID string, user *discordgo.User) bool {
	adminCode := "memebot admin"
	member := getUserMemberFromGuild(discord, guildID, *user)
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
