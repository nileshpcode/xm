package main

import (
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"syscall"
	"xm/app"
	"xm/client"
	"xm/controller"
	"xm/model"
	"xm/repository"
)

func main() {
	xmApp := app.New("XM", app.Config{APIPort: "8080", LogLevel: zerolog.DebugLevel})

	xmApp.DB.AutoMigrate(&model.Company{})

	// initialize app (initializing everything at start to inject dependency)
	xmApp.Initialize(getRoutes(xmApp))

	// run server in a goroutine so that it doesn't block.
	go xmApp.Start()

	xmApp.Logger.Info().Msg("XM service started successfully")

	exitSignal := make(chan os.Signal, 1)
	<-exitSignal

	signal.Notify(exitSignal, os.Interrupt, syscall.SIGTERM)

	// shutdown the API server, waiting for any outstanding requests to complete
	xmApp.Stop()

	xmApp.Logger.Info().Msg("graceful server shutdown complete, exiting")
	os.Exit(0)
}

func getRoutes(xmApp *app.App) []app.RouteSpecifier {
	ipLocationClient := client.NewIpLocationClient("https://ipapi.co")
	companyRepository := repository.NewRepository()
	return []app.RouteSpecifier{controller.NewCompanyController(xmApp, ipLocationClient, companyRepository)}
}
