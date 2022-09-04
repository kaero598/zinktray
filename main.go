package main

import (
	"context"
	"go-fake-smtp/app"
)

func main() {
	storage := &app.Storage{
		Backend: &app.MemoryStorageBackend{},
	}

	smtpServer := app.NewSmtpServer(storage)
	webServer := app.NewWebServer(storage)

	app := app.NewApp(smtpServer, webServer)

	app.Start(context.Background())
}
