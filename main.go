package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/julienschmidt/httprouter"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type Application struct {
	DbClient *bun.DB
	Settings *Settings
}

func (r *Application) startServer(address string) *http.Server {
	router := httprouter.New()
	server := &http.Server{Addr: address, Handler: router}

	// add all handlers
	router.Handle("POST", "/account/register", r.HandleRegistration)

	// listen and serve
	go func() { server.ListenAndServe() }()

	return server
}

func main() {
	// read the settings and handle its presence accordingly
	log.Println("Fetching settings...")

	settings := ReadSettings()
	if settings == nil {
		return
	}

	// configure our application holder
	application := Application{
		Settings: settings,
	}

	// start the server
	log.Println("Preparing server...")
	server := application.startServer(":3030")

	// connect to the database and handle errors accordingly
	pgConfig, pgConfErr := pgx.ParseConfig(settings.PostgresURL)
	if pgConfErr != nil {
		panic("Postgres URL is invalid.")
	}

	postgresDb := stdlib.OpenDB(*pgConfig)
	if postgresDb == nil {
		panic("Could not open connection to db.")
	}

	dbClient := bun.NewDB(postgresDb, pgdialect.New())
	if dbClient == nil {
		panic("Failed to connect to database.")
	}

	application.DbClient = dbClient
	defer dbClient.Close() // ensure the client is closed on shutdown

	models := []interface{}{
		(*UserToken)(nil),
		(*User)(nil),
	}

	for _, model := range models {
		_, err := dbClient.NewCreateTable().Model(model).IfNotExists().Exec(context.Background())
		if err != nil {
			panic(err)
		}
	}

	// create functions & triggers
	funcTrig := []string{
		CreateTokenExpireFunction,
		CreateDeleteTokensOnUserDeletionFunction,
		CreateDeleteTokensOnUserDeletionTrigger,
	}

	for _, query := range funcTrig {
		dbClient.Exec(query)
	}

	// signaling for shutdown
	shutdown := make(chan os.Signal)

	// token ticker
	tokenTicker := time.NewTicker(10 * time.Minute)

	go func() {
		for {
			select {
			case <-tokenTicker.C:
				dbClient.Exec(InvokeTokenExpireFunction)
			case <-shutdown:
				tokenTicker.Stop()
				return
			}
		}
	}()

	// ensure the server termination is handled accordingly
	defer func() {
		if serverErr := server.Close(); serverErr != nil {
			log.Println(serverErr.Error())
		}
	}()

	// ensure an awaited channel
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, os.Interrupt)
	<-shutdown

	// termination message
	log.Println("Cleaning up and shutting down...")
}
