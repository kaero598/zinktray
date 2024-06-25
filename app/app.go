package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"zinktray/app/api"
	"zinktray/app/smtp"
)

// Application represents the application itself.
//
// This slack just delegates the job to its subsystems.
type Application struct {
	// Configured HTTP server
	apiServer *api.Server

	// Configured SMTP server
	smtpServer *smtp.SmtpServer
}

// Start starts all application subsystems and awaits their termination.
func (app *Application) Start(ctx context.Context) {
	appContext, cancel := context.WithCancel(ctx)

	app.watchTerminationSignal(cancel)

	waitGroup := &sync.WaitGroup{}

	waitGroup.Add(1)

	go app.smtpServer.Start(appContext, waitGroup)

	waitGroup.Add(1)

	go app.apiServer.Start(appContext, waitGroup)

	waitGroup.Wait()
}

// watchTerminationSignal wires up OS signal handler and triggers application to terminate
func (app *Application) watchTerminationSignal(cancel context.CancelFunc) {
	channel := make(chan os.Signal, 1)

	signal.Notify(channel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-channel

		cancel()
	}()
}

// NewApp creates new application structure.
func NewApp(smtpServer *smtp.SmtpServer, apiServer *api.Server) *Application {
	return &Application{
		apiServer:  apiServer,
		smtpServer: smtpServer,
	}
}
