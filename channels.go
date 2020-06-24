package main

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
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
	nsfw, name, _, err := GetChannelFromDB(channelID)

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
			go AddChannelToDB(channel.ID, channel.NSFW, channel.Name, channel.GuildID)
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
			go AddChannelToDB(channel.ID, channel.NSFW, channel.Name, channel.GuildID)
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
