package app

import (
	"fmt"
	"github.com/rs/zerolog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"math/rand"
	"os"
	"time"
	"xm/log"
)

// TestApp Provides convenience methods for test
type TestApp struct {
	Application             *App
	controllerRouteProvider func(*App) []RouteSpecifier
	dbInitializer           func(db *gorm.DB)
}

func NewTestApp(name string, controllerRouteProvider func(*App) []RouteSpecifier, dbInitializer func(db *gorm.DB)) *TestApp {
	dbFile := "./test.db?cache=shared&_busy_timeout=60000"

	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())
	randomAPIPort := fmt.Sprintf("10%v%v%v", rand.Intn(9), rand.Intn(9), rand.Intn(9)) // Generating random API port so that if multiple tests can run parallel

	app := &App{name: name, config: Config{APIPort: randomAPIPort}, DB: db}
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	app.Logger = log.New(name, zerolog.DebugLevel, consoleWriter)
	return &TestApp{Application: app, controllerRouteProvider: controllerRouteProvider, dbInitializer: dbInitializer}
}

// Initialize prepares the app for testing
func (testApp *TestApp) Initialize() {
	testApp.Application.Initialize(testApp.controllerRouteProvider(testApp.Application))
	go testApp.Application.Start()
}

// Stop the app
func (testApp *TestApp) Stop() {
	testApp.Application.Stop()
	sqlDB, err := testApp.Application.DB.DB()
	if err != nil {
		sqlDB.Close()
	}
	os.Remove("./test.db")
}

// PrepareEmptyTables clears all table of data
func (testApp *TestApp) PrepareEmptyTables() {
	testApp.dbInitializer(testApp.Application.DB)
}
