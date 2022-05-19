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
	"github.com/imoliwer/sound-point-twitch-bot/server/model"
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

	// configure our application holder
	application := app.Application{
		Settings: settings,
	}

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

	// set up the irc
	twitchIRC := twitch_irc.NewClient()
	twitchIRC.Listen(&application)
	twitchIRC.WithHandler("message", func(client *twitch_irc.Client, state *twitch_irc.TwitchMessageState) {
		log.Println("CHAT:", state)
	})
	twitchIRC.WithHandler("notice", func(client *twitch_irc.Client, state *twitch_irc.TwitchMessageState) {
		log.Println("NOTICE:", state)
	})
	twitchIRC.Join("imoliwer") // TODO: remove before production

	// signal for shutdown
	shutdown := make(chan os.Signal)

	// ensure an awaited channel
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, os.Interrupt)
	<-shutdown

	// close the irc connection
	twitchIRC.Stop()

	// termination message
	log.Println("Cleaning up and shutting down...")
}
