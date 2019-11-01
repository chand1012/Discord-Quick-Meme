// For all functions that manipulate strings
package main

import (
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func stringInSlice(s string, a []string) bool {
	for _, thing := range a {
		if thing == s {
			return true
		}
	}
	return false
}

func textFilter(input string) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9_]+")
	outputString := reg.ReplaceAllString(input, "")
	return outputString
}

// ContainsAnySubstring Checks if any of the substrings in the array are in the test string
func ContainsAnySubstring(testString string, strArray []string) bool {
	for _, str := range strArray {
		if strings.Contains(testString, str) {
			return true
		}
	}
	return false
}

func textFilterSlice(input []string) []string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9_]+")
	var returnSlice []string
	for _, thing := range input {
		output := reg.ReplaceAllString(thing, "")
		returnSlice = append(returnSlice, output)
	}
	return returnSlice
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
