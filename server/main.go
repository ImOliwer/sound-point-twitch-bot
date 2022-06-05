package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imoliwer/sound-point-twitch-bot/server/app"
	"github.com/imoliwer/sound-point-twitch-bot/server/command"
	"github.com/imoliwer/sound-point-twitch-bot/server/model"
	"github.com/imoliwer/sound-point-twitch-bot/server/request"
	"github.com/imoliwer/sound-point-twitch-bot/server/scheduler"
	"github.com/imoliwer/sound-point-twitch-bot/server/sound"
	"github.com/imoliwer/sound-point-twitch-bot/server/twitch_irc"
	"github.com/imoliwer/sound-point-twitch-bot/server/util"
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

	// handle the assignment of request profiles
	twitchAccessory := settings.TwitchAccessory
	request.Profiles = &request.RequestProfiles{
		Twitch: request.TwitchRequestProfile{
			ClientID:     twitchAccessory.ClientID,
			ClientSecret: twitchAccessory.ClientSecret,
			OAuthToken:   twitchAccessory.AuthToken,
			RefreshToken: twitchAccessory.RefreshToken,
		},
	}
	settings.TwitchAccessory = nil // after request assigning

	// handle the validation of the user's Twitch oauth token
	refreshTask := &scheduler.LaterTask{}
	checkToken(true, false, nil, refreshTask)

	validationTask := scheduler.Every(time.Hour, func(_ *scheduler.RepeatingTask) {
		checkToken(false, false, refreshTask, refreshTask)
	})

	defer refreshTask.Cancel()
	defer validationTask.Cancel()

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

	deploymentCover := sound.NewCover(0, 2048)
	{ // set up the main server
		gin.SetMode(gin.ReleaseMode)
		engine := sound.WithCORSAndRecovery(gin.New())

		deploymentCover.Handler(engine)
		sound.RegisterAll(engine, &application, deploymentCover)

		server := &http.Server{
			Addr:    ":9999",
			Handler: engine,
		}

		go func() {
			if err := server.ListenAndServe(); err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					return
				}
				panic(err)
			}
		}()

		util.Log("Deployment & Dashboard", "Server started.")
		defer server.Close()
	}

	{ // set up the twitch irc
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
			panic("Command prefix must consist of ONE character.")
		}

		twitchCmdRegistry := command.NewRegistry(
			twitchCmdPrefix[0],
			map[string]command.PrimaryCommand{
				"points": command.NewPointsCommand(),
			},
			map[string]command.PlaceholderFunc{}, // TODO: register a bunch of general placeholders
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
	application.Settings.Save()
	log.Println("Cleaning up and shutting down...")
}

func checkToken(first bool, ignoreValidation bool, old *scheduler.LaterTask, ptr *scheduler.LaterTask) {
	profile := request.Profiles.Twitch

	if !ignoreValidation {
		if validation := request.ValidateTwitchToken(profile.OAuthToken); validation != nil {
			if first {
				tokenTimer(validation.ExpiresIn, ptr)
			}
			return
		}
	}

	response := request.RefreshCurrentTwitchToken()
	if response == nil {
		panic("Twitch 'Refresh Token' is invalid. Please generate a new one and replace the one in \"settings.json.\"")
	}

	// attempt to revoke current just to avoid multiple (twitch holds up to 50 access tokens per refresh token)
	request.RevokeTwitchToken(request.Profiles.Twitch.OAuthToken)

	request.Profiles.Twitch = request.TwitchRequestProfile{
		ClientID:     profile.ClientID,
		ClientSecret: profile.ClientSecret,
		OAuthToken:   response.AccessToken,
		RefreshToken: response.RefreshToken,
	}

	if old != nil {
		old.Cancel()
	}
	tokenTimer(response.ExpiresIn, ptr)
}

func tokenTimer(expiresIn uint64, ptr *scheduler.LaterTask) {
	*ptr = *scheduler.After(
		time.Duration(expiresIn)*time.Second,
		func(this *scheduler.LaterTask) {
			checkToken(false, true, this, ptr)
		},
	)
}
