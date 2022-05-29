package command

import (
	"errors"
	"strings"

	"github.com/imoliwer/sound-point-twitch-bot/server/request"
	"github.com/imoliwer/sound-point-twitch-bot/server/util"
)

func check_user_and_amount(ctx *Context) bool {
	messages := ctx.AppMessages()
	switch len(ctx.Arguments) {
	case 0:
		ctx.Reply(messages.SpecifyUser)
		return false
	case 1:
		ctx.Reply(messages.SpecifyAmount)
		return false
	}
	return true
}

func amount_of(ctx *Context, allowZero bool) (uint64, error) {
	amount, err := util.Uint64(ctx.Arguments[1])

	if err != nil || (!allowZero && amount == 0) {
		ctx.Reply(ctx.AppMessages().MustSpecifyValidAmount)
		return 0, errors.New("invalid amount")
	}
	return amount, nil
}

func user_of(ctx *Context) (string, uint64, error) {
	var userId uint64
	username := strings.ToLower(ctx.Arguments[0])

	if username == strings.ToLower(ctx.State.User.DisplayName) {
		userId, _ = util.Uint64(ctx.State.User.Id)
	} else {
		userBy := request.TwitchUserBy(username)
		if userBy == nil {
			ctx.Reply(ctx.AppMessages().CouldNotFindUser)
			return username, 0, errors.New("user not found")
		}
		userId, _ = util.Uint64(userBy.Id)
	}
	return username, userId, nil
}
