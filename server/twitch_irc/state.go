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
	"mod":        require_bool_handler,
	"subscriber": require_bool_handler,
	"turbo":      require_bool_handler,
	"user-type":  user_type_handler,
	"badge-info": badge_info_handler,
	"badges":     badges_handler,
	"emotes":     emotes_handler,
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
