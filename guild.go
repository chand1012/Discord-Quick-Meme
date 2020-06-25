package main

import "github.com/bwmarrin/discordgo"

type guildSettings struct {
	Supporter bool
	Proxy     bool
	ProxyMode int8
}

func getServerSettings(discord *discordgo.Session, guildID string) (guildSettings, error) {
	if value, ok := SettingsMap[guildID]; ok {
		return value, nil
	}

	boosted, proxy, proxyMode, err := GetGuildStatus(guildID)

	if err != nil {
		return guildSettings{}, err
	}

	settings := guildSettings{
		Supporter: boosted,
		Proxy:     proxy,
		ProxyMode: proxyMode,
	}

	SettingsMap[guildID] = settings

	return settings, nil
}
