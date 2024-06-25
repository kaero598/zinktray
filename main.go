package main

import (
	"context"
	"zinktray/app"
	"zinktray/app/api"
	"zinktray/app/smtp"
	"zinktray/app/storage"
)

func main() {
	store := storage.NewStorage()

	smtpServer := smtp.NewServer(store)
	apiServer := api.NewServer(store)

	application := app.NewApp(smtpServer, apiServer)

	application.Start(context.Background())
}
