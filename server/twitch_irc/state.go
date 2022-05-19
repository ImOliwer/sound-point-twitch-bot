package twitch_irc

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/imoliwer/sound-point-twitch-bot/server/util"
)

type TwitchNoticeType = uint8
type TwitchSubscriptionTier = uint8
type TwitchUserType = uint8

const (
	NOTICE_SUB TwitchNoticeType = iota + 1
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
	TIER_PRIME TwitchSubscriptionTier = iota + 1
	TIER_1
	TIER_2
	TIER_3
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
	"mod":                         require_bool_handler,
	"subscriber":                  require_bool_handler,
	"turbo":                       require_bool_handler,
	"first-msg":                   require_bool_handler,
	"user-type":                   user_type_handler,
	"badge-info":                  badge_info_handler,
	"badges":                      badges_handler,
	"emotes":                      emotes_handler,
	"msg-id":                      notice_type_handler,
	"msg-param-cumulative-months": require_uint16,
	"msg-param-streak-months":     require_uint16,
	"msg-param-months":            require_uint16,
	"msg-param-sub-plan":          require_uint8,
}

type TwitchUserState struct {
	BadgeInfo    TwitchBadgeInformation `json:"badge_info"`
	Badges       TwitchBadgeList        `json:"badges"`
	Id           string                 `json:"user-id"`
	NameHexColor string                 `json:"color"`
	DisplayName  string                 `json:"display-name"`
	IsModerator  bool                   `json:"mod"`
	IsSubscriber bool                   `json:"subscriber"`
	IsTurbo      bool                   `json:"turbo"`
	Type         TwitchUserType         `json:"user-type"`
}

type TwitchReplyState struct {
	Id              string `json:"reply-parent-msg-id"`
	UserId          string `json:"reply-parent-user-id"`
	UserLogin       string `json:"reply-parent-user-login"`
	UserDisplayName string `json:"reply-parent-display-name"`
	Text            string `json:"reply-parent-msg-body"`
}

type TwitchMessageState struct {
	User           TwitchUserState   `json:"user_state" twitchObj:"true"`
	Reply          TwitchReplyState  `json:"reply_state" twitchObj:"true"`
	Notice         TwitchNoticeState `json:"notice_state" twitchObj:"true"`
	Id             uuid.UUID         `json:"id"`
	ChannelId      string            `json:"room-id"`
	ChannelName    string
	Text           string
	IsFirstMessage bool          `json:"first-msg"`
	Emotes         []TwitchEmote `json:"emotes"`
	BitsCheered    uint32        `json:"bits"`
	ReceivedAt     time.Time     `json:"tmi-sent-ts"`
}

type TwitchSubscriptionState struct {
	CumulativeMonths uint16                 `json:"msg-param-cumulative-months"`
	ShareStreak      bool                   `json:"msg-param-should-share-streak"`
	StreakMonths     uint16                 `json:"msg-param-streak-months"`
	Tier             TwitchSubscriptionTier `json:"msg-param-sub-plan"`
	TierName         string                 `json:"msg-param-sub-plan-name"`
}

type TwitchSubGiftState struct {
	Months      uint16                 `json:"msg-param-months"`
	UserName    string                 `json:"msg-param-recipient-user-name"`
	DisplayName string                 `json:"msg-param-recipient-display-name"`
	Id          string                 `json:"msg-param-recipient-id"`
	Tier        TwitchSubscriptionTier `json:"msg-param-sub-plan"`
	TierName    string                 `json:"msg-param-sub-plan-name"`
	GiftMonths  uint16                 `json:"msg-param-gift-months"`
}

type TwitchNoticeState struct {
	SystemMessage    string                  `json:"system-mg"`
	Login            string                  `json:"login"`
	Subscription     TwitchSubscriptionState `json:"subscription" twitchObj:"true"`
	SubscriptionGift TwitchSubGiftState      `json:"subscription_gift" twitchObj:"true"`
	Type             TwitchNoticeType        `json:"msg-id"`
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

func (r *TwitchMessageState) IsNotice() bool {
	return r.Notice.Type > 0
}

func ProcessMessageState(data []string, t string) TwitchMessageState {
	messageState := TwitchMessageState{
		ChannelName: state_channel_name(data, t),
		Text:        state_text(data, t),
	}
	objectify_irc(data[1], &messageState, objectify_handlers)
	return messageState
}

func state_channel_name(data []string, t string) string {
	if t == "PRIVMSG" {
		return data[5]
	}
	return data[2] // USERNOTICE
}

func state_text(data []string, t string) string {
	if t == "PRIVMSG" {
		return data[6]
	}
	return data[4] // USERNOTICE
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
