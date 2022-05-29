package command

import (
	"context"
	"database/sql"
	"strings"

	"github.com/imoliwer/sound-point-twitch-bot/server/app"
	"github.com/imoliwer/sound-point-twitch-bot/server/model"
	"github.com/imoliwer/sound-point-twitch-bot/server/twitch_irc"
	"github.com/imoliwer/sound-point-twitch-bot/server/util"
)

type UserRequirement = func(*twitch_irc.Client, *twitch_irc.UserState) bool

var ModRequirement UserRequirement = func(c *twitch_irc.Client, us *twitch_irc.UserState) bool {
	return us.IsModerator || us.Badges.Is(twitch_irc.BadgeBroadcaster)
}

type Context struct {
	registry  *Registry
	Client    *twitch_irc.Client
	State     *twitch_irc.MessageState
	Arguments []string
	Temp      map[string]any
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
	commands     map[string]PrimaryCommand
	readyForNext bool
	placeholders map[string]PlaceholderFunc
	Prefix       rune
}

func NewRegistry(prefix rune, initialCmds map[string]PrimaryCommand, placeholders map[string]PlaceholderFunc) Registry {
	if prefix == ' ' {
		prefix = '!'
	}
	registry := Registry{
		Prefix:       prefix,
		commands:     make(map[string]PrimaryCommand),
		readyForNext: true,
		placeholders: placeholders,
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

func (r *Registry) MakeNotReady() {
	r.readyForNext = false
}

func (r *Registry) MakeReady() {
	r.readyForNext = true
}

func (r *Registry) IsReadyForNext() bool {
	return r.readyForNext
}

func (r *Registry) DefaultHandler(client *twitch_irc.Client, state *twitch_irc.MessageState) {
	raw := state.Text
	if raw == "" || len(raw) == 1 || raw[0] != byte(r.Prefix) || !r.readyForNext { // ensure the client is ready for the next command as well
		return
	}

	r.MakeNotReady() // no longer accept commands

	raw = strings.Join(strings.Fields(raw), " ")
	raw = strings.TrimSuffix(raw[1:], "\r")
	arguments := strings.Split(raw, " ")
	name := strings.ToLower(arguments[0])

	command, ok := r.commands[name]
	if !ok {
		r.MakeReady() // accept commands
		return
	}

	primaryOption, ok := client.App.Settings.TwitchBot.Command.Options[name]
	if !try_requirements(command.Command, client, state) || !ok || !primaryOption.Enabled {
		r.MakeReady() // accept commands
		return
	}

	children := command.Children
	arguments = arguments[1:]

	if len(children) > 0 && len(arguments) > 0 {
		childName := strings.ToLower(arguments[0])
		childCommand, ok := children[childName]

		if !ok {
			r.MakeReady() // accept commands
			return
		}

		childOption, ok := primaryOption.Arguments[childName]
		if !try_requirements(childCommand, client, state) || !ok || !childOption.Enabled {
			r.MakeReady() // accept commands
			return
		}

		try_exec(childCommand.Execute, Context{
			Client:    client,
			State:     state,
			Arguments: arguments[1:],
			registry:  r,
			Temp:      make(map[string]any),
		})

		r.MakeReady() // accept commands
		return
	}

	try_exec(command.Execute, Context{
		Client:    client,
		State:     state,
		Arguments: arguments,
		registry:  r,
		Temp:      make(map[string]any),
	})

	r.MakeReady() // accept commands
}

func (r Context) Reply(message string) {
	r.ReplyExtra(message, nil)
}

func (r Context) ReplyExtra(message string, specificPlaceholders map[string]PlaceholderFunc) {
	placeholders := r.registry.placeholders
	if specificPlaceholders != nil {
		placeholders = util.MergeMaps(placeholders, specificPlaceholders)
	}

	r.Client.ReplyTo(
		r.State.Id,
		r.State.ChannelName,
		r.process_placeholders(message, placeholders),
	)
}

func (r Context) CheckErr(err error) bool {
	if err == nil {
		return true
	}

	r.Reply("An error occurred during the operation. Contact personel for further assistance in the matter.")
	return false
}

func (r Context) ModifyUser(mod *model.User, setters map[string][]interface{}) (sql.Result, error) {
	query := r.Client.App.Database.
		NewInsert().
		Model(mod).
		On("CONFLICT (id) DO UPDATE")

	for key, value := range setters {
		query = query.Set(key, value...)
	}
	return query.Exec(context.Background())
}

func (r Context) AppMessages() *app.TwitchCommandMessages {
	return &r.Client.App.Settings.TwitchBot.Command.Messages
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
