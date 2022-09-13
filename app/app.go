package app

import (
	"context"
	"go-fake-smtp/app/smtp"
	"go-fake-smtp/app/web"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// The application. This slack just delegates the job to it's subsystems.
type Application struct {
	// Configured SMTP server
	smtpServer *smtp.SmtpServer

	// Configured HTTP server
	webServer *web.WebServer
}

// Starts all subsystems and awaits their termination
func (app *Application) Start(ctx context.Context) {
	appContext, cancel := context.WithCancel(ctx)

	app.watchTerminationSignal(cancel)

	waitGroup := &sync.WaitGroup{}

	waitGroup.Add(1)

	go app.smtpServer.Start(appContext, waitGroup)

	waitGroup.Add(1)

	go app.webServer.Start(appContext, waitGroup)

	waitGroup.Wait()
}

// Wires up OS signal handler and triggers application to terminate
func (app *Application) watchTerminationSignal(cancel context.CancelFunc) {
	channel := make(chan os.Signal, 1)

	signal.Notify(channel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-channel

		cancel()
	}()
}

// Creates new application
func NewApp(smtpServer *smtp.SmtpServer, webServer *web.WebServer) *Application {
	return &Application{
		smtpServer: smtpServer,
		webServer:  webServer,
	}
}
