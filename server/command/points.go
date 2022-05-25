package command

import (
	"context"
	"strings"

	"github.com/imoliwer/sound-point-twitch-bot/server/model"
	"github.com/imoliwer/sound-point-twitch-bot/server/request"
	"github.com/imoliwer/sound-point-twitch-bot/server/util"
)

func points_give(ctx Context) {
	switch len(ctx.Arguments) {
	case 0:
		ctx.Reply("Please specify a user to give points to.")
		return
	case 1:
		ctx.Reply("Please specify the amount of points to give.")
		return
	}

	amount, err := util.Uint64(ctx.Arguments[1])
	if err != nil || amount == 0 {
		ctx.Reply("You must specify a valid amount.")
		return
	}

	var userId uint64
	username := strings.ToLower(ctx.Arguments[0])

	if username == strings.ToLower(ctx.State.User.DisplayName) {
		userId, _ = util.Uint64(ctx.State.User.Id)
	} else {
		userBy := request.TwitchUserBy(username)
		if userBy == nil {
			ctx.Reply("Could not find the user \"%s.\"", username)
			return
		}
		userId, _ = util.Uint64(userBy.Id)
	}

	_, err = ctx.Client.App.Database.
		NewInsert().
		Model(&model.User{
			ID:     userId,
			Points: amount,
		}).
		On("CONFLICT (id) DO UPDATE").
		Set("points = points + ?", amount).
		Exec(context.Background())

	if err != nil {
		ctx.Reply("An error occurred during the operation. Contact personel for further assistance in the matter.")
		return
	}

	ctx.Reply("\"%s\" has been given %d points.", username, amount)
}

func NewPointsCommand() PrimaryCommand {
	return PrimaryCommand{
		Command: Command{
			Requirements: []UserRequirement{
				ModRequirement,
			},
		},
		Children: map[string]Command{
			"give": {
				Requirements: make([]UserRequirement, 0),
				Execute:      points_give,
			},
		},
	}
}
