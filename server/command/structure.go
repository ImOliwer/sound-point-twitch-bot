package command

import (
	"strings"

	"github.com/imoliwer/sound-point-twitch-bot/server/twitch_irc"
	"github.com/imoliwer/sound-point-twitch-bot/server/util"
)

type Requirement = func(*twitch_irc.UserState) bool

type Context struct {
	State     *twitch_irc.MessageState
	Arguments []string
}

type Command struct {
	Children     map[string]Command
	Requirements []Requirement
	Execution    func(*twitch_irc.Client, *twitch_irc.MessageState)
}

type Registry struct {
	parents map[string]Command
	Prefix  rune
}

func (r *Registry) Include(name string, parent Command) {
	lowered := strings.ToLower(name)
	r.parents[lowered] = parent
	util.Log("Commands", "Registered '%s'", lowered)
}

func (r *Registry) Exclude(name string) {
	lowered := strings.ToLower(name)
	if _, ok := r.parents[lowered]; ok {
		delete(r.parents, lowered)
		util.Log("Commands", "Unregistered '%s'", lowered)
	}
}
