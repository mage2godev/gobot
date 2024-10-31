package main

import (
	"context"
	"flag"
	"log"

	tgClient "GoLangProjects/clients/telegram"
	event_consumer "GoLangProjects/consumer/event-consumer"
	"GoLangProjects/events/telegram"
	"GoLangProjects/storage/sqlite"
)

const (
	tgBotHost         = "api.telegram.org"
	sqliteStoragePath = "data/sqlite/storage.db"
	batchSize         = 100
)

// -tg-bot-token '7773342135:AAG9ozhWHbZnkdoPCxJOLe34wW5Elgt2w_A'

func main() {

	//s := files.New(storagePath)
	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatal("can't connect to storage: ", err)
	}

	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can't init storage: ", err)
	}

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		s,
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access telegram",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("toke is not exist")
	}

	return *token
}
