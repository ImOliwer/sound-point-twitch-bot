package twitch_irc

import (
	"strings"

	"github.com/imoliwer/sound-point-twitch-bot/server/util"
)

// transform 0->1 through false->true
func require_bool_handler(value string) interface{} {
	return util.RequireBool(value)
}

// get user type by string
func user_type_handler(value string) interface{} {
	switch value {
	case "admin":
		return USER_ADMIN
	case "global_mod":
		return USER_GLOBAL_MOD
	case "staff":
		return USER_STAFF
	case "mod":
		return USER_MOD
	default:
		return USER_NORMAL
	}
}

// convert all badge information from raw string
func badge_info_handler(value string) interface{} {
	badgeInformation := TwitchBadgeInformation{
		Subscription:      0,
		SubscriptionGifts: 0,
	}

	for badgeName, badgeValue := range split_raw(value) {
		parsed := util.Uint32(badgeValue)
		switch badgeName {
		case "subscriber":
			badgeInformation.Subscription = parsed
			break
		case "sub-gifter":
			badgeInformation.SubscriptionGifts = parsed
			break
		}
	}
	return badgeInformation
}

// convert all badges in a string
func badges_handler(value string) interface{} {
	badgeList := TwitchBadgeList{single_versions: 0}
	var singleVersions util.Flag

	// handle
	for badgeName, badgeValue := range split_raw(value) {
		switch badgeName {
		case "admin":
			singleVersions.Append(BADGE_ADMIN)
		case "staff":
			singleVersions.Append(BADGE_STAFF)
		case "moderator":
			singleVersions.Append(BADGE_MODERATOR)
		case "partner":
			singleVersions.Append(BADGE_PARTNER)
		case "turbo":
			singleVersions.Append(BADGE_TURBO)
		case "glhf-pledge":
			singleVersions.Append(BADGE_GLHF_PLEDGE)
		case "broadcaster":
			singleVersions.Append(BADGE_BROADCASTER)
		case "bits-charity":
			singleVersions.Append(BADGE_BITS_CHARITY)
		case "subscriber":
			badgeList.Subscriber = util.Uint32(badgeValue)
		case "bits":
			badgeList.Bits = util.Uint32(badgeValue)
		}
	}

	badgeList.single_versions = singleVersions
	return badgeList
}

// convert all the emote positions accordingly from raw value
func emotes_handler(value string) interface{} {
	emotes := make([]TwitchEmote, 0)
	if value != "" {
		for _, raw := range strings.Split(value, "/") {
			slice := strings.Split(raw, ":")
			twitchEmote := TwitchEmote{Id: slice[0]}

			positions := make([]TwitchEmotePosition, 0)
			for _, rawPosition := range strings.Split(slice[1], ",") {
				whole := strings.Split(rawPosition, "-")
				positions = append(positions, TwitchEmotePosition{
					StartPos: util.Uint16(whole[0]),
					EndPos:   util.Uint16(whole[1]),
				})
			}

			twitchEmote.Positions = positions
			emotes = append(emotes, twitchEmote)
		}
	}
	return emotes
}
