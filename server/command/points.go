package command

import (
	"context"

	"github.com/imoliwer/sound-point-twitch-bot/server/model"
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

	message := ctx.AppMessages().PointsNoArg
	if err != nil {
		ctx.Reply(message, 0)
		return
	}

	ctx.Reply(message, response.Points)
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
	ctx.Reply(ctx.AppMessages().PointsGiveSuccess, username, amount)
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
	ctx.Reply(ctx.AppMessages().PointsSetSuccess, username, amount)
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
