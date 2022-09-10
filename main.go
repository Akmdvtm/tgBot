package main

import (
	"flag"
	"log"
	"main.go/client/telegram"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	tgClient := telegram.New(tgBotHost, mustToken()) // инициализировали client, mustToken получает токен
	//fetcher = fetcher.New(tgClient)
	//processor = processor.New(tgClient)
	//
	//consumer.Start(fetcher, processor)
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
