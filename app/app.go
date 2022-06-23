package app

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
	"xm/event"
	"xm/log"
)

// App structure for xm microservice
type App struct {
	name   string
	config Config
	DB     *gorm.DB
	Router *mux.Router
	server *http.Server
	Logger *zerolog.Logger
	event.Dispatcher
}

// Config consists config fields needed to start the app
type Config struct {
	APIPort  string
	LogLevel zerolog.Level
}

func New(name string, config Config) *App {
	app := &App{name: name, config: config}
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	app.Logger = log.New(name, config.LogLevel, consoleWriter)
	app.initializeDB()
	app.initializeQueue()
	return app
}

// Initialize initializes properties of the app
func (app *App) Initialize(routeSpecifiers []RouteSpecifier) {

	logger := app.Logger
	app.Router = mux.NewRouter()
	app.Router.Use(mux.CORSMethodMiddleware(app.Router))

	for _, routeSpecifier := range routeSpecifiers {
		routeSpecifier.RegisterRoutes(app.Router)
	}

	logger.Debug().Str("app", app.name).Msg("Api server will start on port: " + app.config.APIPort)

	app.server = &http.Server{
		Addr:    "0.0.0.0:" + app.config.APIPort,
		Handler: app.Router,
	}
}

// initializeDB connects to db
func (app *App) initializeDB() error {
	db, err := gorm.Open(sqlite.Open("xm.db"), &gorm.Config{})
	if err != nil {
		app.Logger.Fatal().Err(err).Msg("failed to connect database, exiting the application!")
	}
	app.DB = db
	return nil
}

// initializeQueue connects to queue
func (app *App) initializeQueue() error {
	dispatcher, err := event.NewRabbitMQEventDispatcher(app.Logger, getQueueConnectionString())
	if err != nil {
		app.Logger.Fatal().Err(err).Msg("failed to connect rabbitmq, exiting the application!")
	}
	app.Dispatcher = dispatcher
	return nil
}

//Start http server and start listening to the requests
func (app *App) Start() {
	if err := app.server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			app.Logger.Fatal().Err(err).Msg("Unable to start server, exiting the application!")
		}
	}
}

// Stop http server
func (app *App) Stop() {
	wait, _ := time.ParseDuration("1m")
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	app.server.Shutdown(ctx)
}

// RouteSpecifier should be implemented by the class that sets routes for the API endpoints
type RouteSpecifier interface {
	RegisterRoutes(router *mux.Router)
}

func getQueueConnectionString() string {
	var queueHost, queuePort, queueUser, queuePassword string

	queueProtocol := "amqp"
	queueHost, ok := os.LookupEnv("QUEUE_HOST")
	if !ok {
		queueHost = "localhost"
	}
	queuePassword, ok = os.LookupEnv("QUEUE_PWD")
	if !ok {
		queuePassword = "guest"
	}
	queueUser, ok = os.LookupEnv("QUEUE_USER")
	if !ok {
		queueUser = "guest"
	}
	queuePort, ok = os.LookupEnv("QUEUE_PORT")
	if !ok {
		queuePort = "5672"
	}

	return fmt.Sprintf("%v://%v:%v@%v:%v/", queueProtocol, queueUser, queuePassword, queueHost, queuePort)
}
