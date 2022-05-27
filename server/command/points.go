package command

import (
	"context"
	"errors"
	"strings"

	"github.com/imoliwer/sound-point-twitch-bot/server/model"
	"github.com/imoliwer/sound-point-twitch-bot/server/request"
	"github.com/imoliwer/sound-point-twitch-bot/server/util"
)

func points_no_arg(ctx Context) {
	userId, _ := util.Uint64(ctx.State.User.Id)
	var response model.User

	err := ctx.Client.App.Database.
		NewSelect().
		Model(&response).
		Where("id = ?", userId).
		Scan(context.Background())

	if err != nil {
		ctx.Reply("You currently have 0 points.")
		return
	}
	ctx.Reply("You currently have %d points.", response.Points)
}

func points_give(ctx Context) {
	amount, username, userId, ok := points_standard(&ctx, false)
	if !ok {
		return
	}

	_, err := ctx.ModifyUser(
		&model.User{ID: userId, Points: amount},
		map[string][]interface{}{
			"points = points + ?": {amount},
		},
	)

	ctx.CheckErr(err)
	ctx.Reply("%s has been given %d points.", username, amount)
}

func points_set(ctx Context) {
	amount, username, userId, ok := points_standard(&ctx, true)
	if !ok {
		return
	}

	_, err := ctx.ModifyUser(
		&model.User{ID: userId, Points: amount},
		map[string][]interface{}{
			"points = ?": {amount},
		},
	)

	ctx.CheckErr(err)
	ctx.Reply("The points of %s has been set to %d.", username, amount)
}

func NewPointsCommand() PrimaryCommand {
	modRequirements := []UserRequirement{
		ModRequirement,
	}
	return PrimaryCommand{
		Command: Command{
			Requirements: make([]UserRequirement, 0),
			Execute:      points_no_arg,
		},
		Children: map[string]Command{
			"give": {
				Requirements: modRequirements,
				Execute:      points_give,
			},
			"set": {
				Requirements: modRequirements,
				Execute:      points_set,
			},
		},
	}
}

//////////////////////
// HELPER FUNCTIONS //
//////////////////////

func points_standard(ctx *Context, allowZero bool) (uint64, string, uint64, bool) {
	if !check_user_and_amount(ctx) {
		return 0, "", 0, false
	}

	amount, err := amount_of(ctx, allowZero)
	if err != nil {
		return 0, "", 0, false
	}

	username, userId, err := user_of(ctx)
	if err != nil {
		return 0, "", 0, false
	}
	return amount, username, userId, true
}

func check_user_and_amount(ctx *Context) bool {
	switch len(ctx.Arguments) {
	case 0:
		ctx.Reply("Please specify a user.")
		return false
	case 1:
		ctx.Reply("Please specify the amount of points.")
		return false
	}
	return true
}

func amount_of(ctx *Context, allowZero bool) (uint64, error) {
	amount, err := util.Uint64(ctx.Arguments[1])
	if err != nil || (!allowZero && amount == 0) {
		ctx.Reply("You must specify a valid amount.")
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
			ctx.Reply("Could not find the user \"%s.\"", username)
			return username, 0, errors.New("user not found")
		}
		userId, _ = util.Uint64(userBy.Id)
	}
	return username, userId, nil
}
