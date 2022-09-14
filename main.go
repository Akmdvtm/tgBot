package main

import (
	"flag"
	"log"
	tgClient "tgBot/client/telegram"
	eventConsumer "tgBot/consumer/event-consumer"
	"tgBot/events/telegram"
	"tgBot/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {
	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Println("servise started")

	consumer := eventConsumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal()
	}
}

func mustToken() string { // must обознач !return err и с ней работать осторожно
	token := flag.String("bot-token",
		"",
		"token for access for tgbot",
	) // 1 имя 2 по дефолту 3 подсказка (писать осмысленно)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified") // стандартная функция Go, при fatal программа аварийно завершается
	}

	return *token
}
