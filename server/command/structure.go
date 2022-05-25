package command

import (
	"strings"

	"github.com/imoliwer/sound-point-twitch-bot/server/twitch_irc"
	"github.com/imoliwer/sound-point-twitch-bot/server/util"
)

type UserRequirement = func(*twitch_irc.Client, *twitch_irc.UserState) bool

var ModRequirement UserRequirement = func(c *twitch_irc.Client, us *twitch_irc.UserState) bool {
	return us.IsModerator || us.Badges.Is(twitch_irc.BadgeBroadcaster)
}

type Context struct {
	Client    *twitch_irc.Client
	State     *twitch_irc.MessageState
	Arguments []string
}

type Command struct {
	Requirements []UserRequirement
	Execute      func(Context)
}

type PrimaryCommand struct {
	Command
	Children map[string]Command
}

type Registry struct {
	commands map[string]PrimaryCommand
	Prefix   rune
}

func NewRegistry(prefix rune, initialCmds map[string]PrimaryCommand) Registry {
	if prefix == ' ' {
		prefix = '!'
	}
	registry := Registry{
		Prefix:   prefix,
		commands: make(map[string]PrimaryCommand),
	}
	for key, cmd := range initialCmds {
		registry.Include(key, cmd) // this is used only for the loggings and lowered names
	}
	return registry
}

func (r *Registry) Include(name string, parent PrimaryCommand) {
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

func (r *Registry) DefaultHandler(client *twitch_irc.Client, state *twitch_irc.MessageState) {
	raw := state.Text
	if raw == "" || len(raw) == 1 || raw[0] != byte(r.Prefix) {
		return
	}

	raw = strings.TrimSuffix(raw[1:], "\r")
	arguments := strings.Split(raw, " ")
	name := strings.ToLower(arguments[0])

	command, ok := r.commands[name]
	if !ok || !try_requirements(command.Command, client, state) {
		return
	}

	children := command.Children
	arguments = arguments[1:]

	if len(children) > 0 {
		if len(arguments) == 0 {
			return
		}

		childCommand, ok := children[strings.ToLower(arguments[0])]
		if !ok || !try_requirements(childCommand, client, state) {
			return
		}

		try_exec(childCommand.Execute, Context{
			Client:    client,
			State:     state,
			Arguments: arguments[1:],
		})
		return
	}

	try_exec(command.Execute, Context{
		Client:    client,
		State:     state,
		Arguments: arguments,
	})
}

func (r Context) Reply(message string, args ...any) {
	r.Client.ReplyTo(r.State.Id, r.State.ChannelName, message, args...)
}

func try_requirements(cmd Command, client *twitch_irc.Client, state *twitch_irc.MessageState) bool {
	for _, requirement := range cmd.Requirements {
		if !requirement(client, &state.User) {
			return false
		}
	}
	return true
}

func try_exec(exec func(Context), ctx Context) {
	if exec != nil {
		exec(ctx)
	}
}
