package main

import (
	"math/rand"
	"strings"
)

const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand

func discordChannelFilter(req string) bool {
	if getDiscordConfigBool("channels.filter") {
		if strings.Contains(getDiscordChannels(), req) {
			return true
		}
		if getDiscordKOMChannel(req) {
			return true
		}
		return false
	}
	return true
}

func discordAuthorRolePermissionCheck(roles []string) (bool, string) {
	for _, x := range roles {
		for _, y := range getDiscordGroupRoles("admin") {
			if x == y {
				return true, "admin"
			}
		}
	}
	return false, ""
}

func discordMessageHandler(dpack DataPackage) {
	superdebug("In discord message handler")
	// If the string doesn't have the prefix parse as text, if it does parse as a command.
	if !strings.HasPrefix(dpack.Message, getDiscordConfigString("prefix")) {
		superdebug("checking keywords")
		dpack.MsgTye = "keyword"
		if discordChannelFilter(dpack.ChannelID) {
			debug("No prefix was found parsing for keywords.")
			parseKeyword(dpack)
		}
	} else {
		superdebug("Checking keywords")
		dpack.MsgTye = "command"
		// if there is a prefix check permissions on the user and run commands per group.
		dpack.Message = strings.TrimPrefix(dpack.Message, getDiscordConfigString("prefix"))
		if dpack.Perms {
			if dpack.Group == "admin" {
				parseAdminCommand(dpack)
				parseModCommand(dpack)
			}
			if dpack.Group == "mod" {
				parseModCommand(dpack)
			}
		}
		// parse commands for matches
		debug("Prefix was found parsing for commands.")
		dpack.Message = strings.TrimPrefix(dpack.Message, getDiscordConfigString("prefix"))
		parseCommand(dpack)
		// remove previous commands if discord.command.remove is true
		if getDiscordConfigBool("command.remove") {
			if getCommandStatus(dpack.Message) {
				deleteDiscordMessage(dpack)
				debug("Cleared command message.")
			}
			if strings.HasPrefix(dpack.Message, "list") || strings.HasPrefix(dpack.Message, "ggl") {
				deleteDiscordMessage(dpack)
				debug("Cleared command message.")
			}
		}
	}
}

func discordImageRandGen() string {
	b := make([]byte, 12)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
