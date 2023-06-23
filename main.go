package main

import (
	"context"
	"go-fake-smtp/app"
	"go-fake-smtp/app/api"
	"go-fake-smtp/app/smtp"
	"go-fake-smtp/app/storage"
)

func main() {
	store := storage.NewStorage()

	smtpServer := smtp.NewServer(store)
	apiServer := api.NewServer(store)

	application := app.NewApp(smtpServer, apiServer)

	application.Start(context.Background())
}
