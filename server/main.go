package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/imoliwer/sound-point-twitch-bot/server/app"
	"github.com/imoliwer/sound-point-twitch-bot/server/command"
	"github.com/imoliwer/sound-point-twitch-bot/server/model"
	"github.com/imoliwer/sound-point-twitch-bot/server/request"
	"github.com/imoliwer/sound-point-twitch-bot/server/twitch_irc"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

func main() {
	// read the settings and handle its presence accordingly
	log.Println("Fetching settings...")

	settings := app.ReadSettings()
	if settings == nil {
		return
	}

	// configure our application holders
	application := app.Application{
		Settings: settings,
	}

	// TODO: validate the oauth token and fetch a new one via refresh token if invalid before proceeding with assigning the Twitch profile
	request.Assign(&request.RequestProfiles{
		Twitch: request.TwitchRequestProfile{
			ClientID:   settings.TwitchAccessory.ClientID,
			OAuthToken: settings.TwitchAccessory.AuthToken,
		},
	})

	settings.TwitchAccessory = nil // after request assigning

	// attempt to create the SQLite database in case it's absent
	if _, err := os.Stat("data.db"); errors.Is(err, os.ErrNotExist) {
		dataFile, dataFileErr := os.Create("data.db")
		if dataFileErr != nil {
			panic("Failed to create database.")
		}
		dataFile.Close()
	}

	// show that we're connecting
	log.Println("Attempting to connect to database...")

	// connect to the database and handle errors accordingly
	sqlDb, sqlDbErr := sql.Open(sqliteshim.ShimName, "./data.db")
	if sqlDbErr != nil {
		panic("Could not open connection to SQLite database.")
	}

	db := bun.NewDB(sqlDb, sqlitedialect.New())
	if db == nil {
		panic("Could not open connection via BUN.")
	}

	log.Println("Connected to database successfully.")

	application.Database = db
	defer db.Close() // ensure the client is closed on shutdown

	// prepare table(s)
	modelArray := []interface{}{
		(*model.User)(nil),
	}

	for _, model := range modelArray {
		_, err := db.NewCreateTable().Model(model).IfNotExists().Exec(context.Background())
		if err != nil {
			panic(err)
		}
	}

	// set up the twitch irc
	{
		twitchChannelToJoin := settings.TwitchBot.Channel
		if twitchChannelToJoin == "" {
			panic("Invalid channel name in settings.")
		}

		twitchIRC := twitch_irc.NewClient(&application)
		twitchIRC.Listen()
		defer twitchIRC.Stop()

		// handle commands
		twitchCmdPrefix := []rune(settings.TwitchBot.Command.Prefix)
		if len(twitchCmdPrefix) != 1 {
			panic("Command prefix must consist of ONE character")
		}

		twitchCmdRegistry := command.NewRegistry(
			twitchCmdPrefix[0],
			map[string]command.PrimaryCommand{
				"points": command.NewPointsCommand(),
			},
		)

		twitchIRC.WithHandler("message", twitchCmdRegistry.DefaultHandler)
		twitchIRC.Join(twitchChannelToJoin) // join after command handle
	}

	// signal for shutdown
	shutdown := make(chan os.Signal)

	// ensure an awaited channel
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, os.Interrupt)
	<-shutdown

	// termination message
	log.Println("Cleaning up and shutting down...")
}
