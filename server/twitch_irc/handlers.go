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
		return TierPrime
	case "1000":
		return Tier1
	case "2000":
		return Tier2
	case "3000":
		return Tier3
	default:
		return uint8(0)
	}
}

// get user type by string
func user_type_handler(value string) interface{} {
	switch value {
	case "admin":
		return UserAdmin
	case "global_mod":
		return UserGlobalMod
	case "staff":
		return UserStaff
	case "mod":
		return UserMod
	default:
		return UserNormal
	}
}

// fetch the exact notice type by value
func notice_type_handler(value string) interface{} {
	switch value {
	case "sub":
		return NoticeSub
	case "resub":
		return NoticeResub
	case "subgift":
		return NoticeSubGift
	case "giftpaidupgrade":
		return NoticeGiftPaidUpgrade
	case "rewardgift":
		return NoticeRewardGift
	case "anongiftpaidupgrade":
		return NoticeAnonGiftPaidUpgrade
	case "raid":
		return NoticeRaid
	case "unraid":
		return NoticeUnraid
	case "bitsbadgetier":
		return NoticeBitsBadgeTier
	default:
		return NoticeRitual
	}
}

// convert all badge information from raw string
func badge_info_handler(value string) interface{} {
	badgeInformation := BadgeInformation{
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
	badgeList := BadgeList{single_versions: 0}
	var singleVersions util.Flag

	// handle
	for badgeName, badgeValue := range split_raw(value) {
		switch badgeName {
		case "admin":
			singleVersions.Append(BadgeAdmin)
		case "staff":
			singleVersions.Append(BadgeStaff)
		case "moderator":
			singleVersions.Append(BadgeModerator)
		case "partner":
			singleVersions.Append(BadgePartner)
		case "turbo":
			singleVersions.Append(BadgeTurbo)
		case "glhf-pledge":
			singleVersions.Append(BadgeGlhfPledge)
		case "broadcaster":
			singleVersions.Append(BadgeBroadcaster)
		case "bits-charity":
			singleVersions.Append(BadgeBitsCharity)
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
	emotes := make([]Emote, 0)
	if value != "" {
		for _, raw := range strings.Split(value, "/") {
			slice := strings.Split(raw, ":")
			twitchEmote := Emote{Id: slice[0]}

			positions := make([]EmotePosition, 0)
			for _, rawPosition := range strings.Split(slice[1], ",") {
				whole := strings.Split(rawPosition, "-")
				positions = append(positions, EmotePosition{
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
