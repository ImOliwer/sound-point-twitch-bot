package app

import (
	"github.com/imoliwer/sound-point-twitch-bot/server/model"
	"github.com/uptrace/bun"
)

type Application struct {
	Database *bun.DB
	Settings *Settings
}

type ModelStructure struct {
	User *model.User
}
