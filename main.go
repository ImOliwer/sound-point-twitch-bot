package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/julienschmidt/httprouter"
)

type Application struct {
	DbClient *pg.DB
	Settings *Settings
}

func (r *Application) startServer(address string) *http.Server {
	router := httprouter.New()
	server := &http.Server{Addr: address, Handler: router}

	// add all handlers
	router.Handle("POST", "/account/register", HandleRegistration)

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
	connectOpts, connectOptsErr := pg.ParseURL(settings.PostgresURL)
	if connectOptsErr != nil {
		panic(connectOptsErr)
	}

	dbClient := pg.Connect(connectOpts)
	if dbClient == nil {
		panic("Failed to connect to database.")
	}

	application.DbClient = dbClient
	defer dbClient.Close() // ensure the client is closed on shutdown

	models := []interface{}{
		(*User)(nil),
	}

	for _, model := range models {
		err := dbClient.Model(model).CreateTable(&orm.CreateTableOptions{IfNotExists: true})
		if err != nil {
			panic(err)
		}
	}

	// ensure the server termination is handled accordingly
	defer func() {
		if serverErr := server.Close(); serverErr != nil {
			log.Println(serverErr.Error())
		}
	}()

	// ensure an awaited channel
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, os.Interrupt)
	<-shutdown

	// termination message
	log.Println("Cleaning up and shutting down...")
}
