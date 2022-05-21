package command

import (
	"strings"

	"github.com/imoliwer/sound-point-twitch-bot/server/twitch_irc"
	"github.com/imoliwer/sound-point-twitch-bot/server/util"
)

type UserRequirement = func(*twitch_irc.Client, *twitch_irc.UserState) bool

type Context struct {
	State     *twitch_irc.MessageState
	Arguments []string
}

type Command struct {
	Requirements []UserRequirement
	Execution    func(*twitch_irc.Client, Context)
}

type Registry struct {
	commands map[string]Command
	Prefix   rune
}

func (r *Registry) Include(name string, parent Command) {
	lowered := strings.ToLower(name)
	r.commands[lowered] = parent
	util.Log("Commands", "Registered '%s'", lowered)
}

func (r *Registry) Exclude(name string) {
	lowered := strings.ToLower(name)
	if _, ok := r.commands[lowered]; ok {
		delete(r.commands, lowered)
		util.Log("Commands", "Unregistered '%s'", lowered)
	}
}

func (r *Registry) Handle(raw string, client *twitch_irc.Client, state *twitch_irc.MessageState) {
	if raw == "" || len(raw) == 1 || raw[0] != byte(r.Prefix) {
		return
	}

	arguments := strings.Split(raw[1:], " ")
	name := strings.ToLower(arguments[0])

	command, ok := r.commands[name]
	if !ok {
		return
	}

	for _, requirement := range command.Requirements {
		if !requirement(client, &state.User) {
			return
		}
	}

	command.Execution(
		client,
		Context{
			State:     state,
			Arguments: arguments[1:],
		},
	)
}
