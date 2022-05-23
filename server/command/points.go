package command

import "github.com/imoliwer/sound-point-twitch-bot/server/twitch_irc"

func points_give_handle(c *twitch_irc.Client, ctx Context) {
	c.Chat(ctx.State.ChannelName, "callback") // FIXME: just for testing purposes
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
				Execute:      points_give_handle,
			},
		},
	}
}
