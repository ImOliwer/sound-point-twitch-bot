package twitch_irc

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/imoliwer/sound-point-twitch-bot/server/util"
)

type TwitchNoticeType = uint8
type TwitchUserType = uint8

const (
	NOTICE_SUB TwitchNoticeType = iota
	NOTICE_SUB_GIFT
	NOTICE_RESUB
	NOTICE_GIFT_PAID_UPGRADE
	NOTICE_ANON_GIFT_PAID_UPGRADE
	NOTICE_REWARD_GIFT
	NOTICE_RAID
	NOTICE_UNRAID
	NOTICE_RITUAL
	NOTICE_BITS_BADGE_TIER
)

const (
	USER_NORMAL TwitchUserType = iota
	USER_MOD
	USER_ADMIN
	USER_GLOBAL_MOD
	USER_STAFF
)

const (
	BADGE_ADMIN util.Flag = 1 << iota
	BADGE_STAFF
	BADGE_MODERATOR
	BADGE_PARTNER
	BADGE_TURBO
	BADGE_GLHF_PLEDGE
	BADGE_BROADCASTER
	BADGE_BITS_CHARITY
)

var objectify_handlers = map[string]objectify_handler{
	"id": func(value string) interface{} {
		return uuid.MustParse(value)
	},
	"bits": func(value string) interface{} {
		return util.Uint32(value)
	},
	"tmi-sent-ts": func(value string) interface{} {
		return time.UnixMilli(util.Int64(value))
	},
	"mod": func(value string) interface{} {
		return util.RequireBool(value)
	},
	"subscriber": func(value string) interface{} {
		return util.RequireBool(value)
	},
	"turbo": func(value string) interface{} {
		return util.RequireBool(value)
	},
	"user-type": func(value string) interface{} {
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
	},
	"badge-info": func(value string) interface{} {
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
	},
	"badges": func(value string) interface{} {
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
	},
	"emotes": func(value string) interface{} {
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
	},
}

type TwitchUserState struct {
	BadgeInfo    TwitchBadgeInformation
	Badges       TwitchBadgeList
	Id           string
	NameHexColor string
	DisplayName  string
	IsModerator  bool
	IsSubscriber bool
	IsTurbo      bool
	Type         TwitchUserType
}

type TwitchMessageState struct {
	User        TwitchUserState `json:"user_state" link:"badge-info=BadgeInfo;badges=Badges;user-id=Id;color=NameHexColor;display-name=DisplayName;mod=IsModerator;subscriber=IsSubscriber;turbo=IsTurbo;user-type=Type"`
	Notice      TwitchNoticeState
	Id          uuid.UUID `json:"id"`
	ChannelId   string    `json:"room-id"`
	ChannelName string
	Text        string
	Emotes      []TwitchEmote `json:"emotes"`
	BitsCheered uint32        `json:"bits"`
	ReceivedAt  time.Time     `json:"tmi-sent-ts"`
}

type TwitchNoticeState struct {
	SystemMessage  string
	MessageId      string
	LoginName      string
	TargetUserId   string
	IsFirstMessage bool
	Type           TwitchNoticeType
}

type TwitchBadgeInformation struct {
	Subscription      uint32
	SubscriptionGifts uint32
}

type TwitchBadgeList struct {
	// single-version badges
	single_versions util.Flag
	// multiple-version badges; those of which are not present for said user, result in a "default" value of -1.
	Subscriber uint32
	Bits       uint32
}

type TwitchEmote struct {
	Id        string
	Positions []TwitchEmotePosition
}

type TwitchEmotePosition struct {
	StartPos uint16
	EndPos   uint16
}

func (r *TwitchBadgeList) Is(flag util.Flag) bool {
	return r.single_versions.Has(flag)
}

func ProcessMessageState(data []string) TwitchMessageState {
	messageState := TwitchMessageState{
		ChannelName: data[5],
		Text:        data[6],
	}
	objectify_irc(data[1], &messageState, objectify_handlers)
	return messageState
}

func split_raw(value string) map[string]string {
	population := map[string]string{}
	for _, raw := range strings.Split(value, ",") {
		if raw == "" {
			continue
		}

		split := strings.Split(raw, "/")
		if len(split) == 0 {
			continue
		}

		key := split[0]
		population[key] = split[1]
	}
	return population
}
