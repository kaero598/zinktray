package main

import (
	"context"
	"go-fake-smtp/app"
	"go-fake-smtp/app/smtp"
	"go-fake-smtp/app/storage"
	"go-fake-smtp/app/web"
)

func main() {
	storage := storage.NewStorage(&storage.MemoryBackend{})

	smtpServer := smtp.NewServer(storage)
	webServer := web.NewServer(storage)

	app := app.NewApp(smtpServer, webServer)

	app.Start(context.Background())
}
