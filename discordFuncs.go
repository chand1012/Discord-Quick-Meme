// for any misc discord functions
package main

import (
	"fmt"
	"strconv"
	"sync"

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
