package twitch_irc

import (
	"strings"

	"github.com/imoliwer/sound-point-twitch-bot/server/util"
)

// transform 0->1 through false->true
func bool_handler(value string) interface{} {
	if value == "" {
		return false
	}
	return util.RequireBool(value)
}

// convert string into uint16
func uint16_handler(value string) interface{} {
	if value == "" {
		return 0
	}
	return util.Uint16(value)
}

// convert string into a sub plan
func sub_plan_handler(value string) interface{} {
	switch value {
	case "Prime":
		return TIER_PRIME
	case "1000":
		return TIER_1
	case "2000":
		return TIER_2
	case "3000":
		return TIER_3
	default:
		return uint8(0)
	}
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

// fetch the exact notice type by value
func notice_type_handler(value string) interface{} {
	switch value {
	case "sub":
		return NOTICE_SUB
	case "resub":
		return NOTICE_RESUB
	case "subgift":
		return NOTICE_SUB_GIFT
	case "giftpaidupgrade":
		return NOTICE_GIFT_PAID_UPGRADE
	case "rewardgift":
		return NOTICE_REWARD_GIFT
	case "anongiftpaidupgrade":
		return NOTICE_ANON_GIFT_PAID_UPGRADE
	case "raid":
		return NOTICE_RAID
	case "unraid":
		return NOTICE_UNRAID
	case "bitsbadgetier":
		return NOTICE_BITS_BADGE_TIER
	default:
		return NOTICE_RITUAL
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
