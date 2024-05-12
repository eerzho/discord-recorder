package main

import (
	"log"

	"discord-recorder/internal/app/discord_bot"
	"discord-recorder/internal/app_log"
	"discord-recorder/internal/app_minio"
	"discord-recorder/internal/config"
)

func main() {
	log.Print("parsing configuration")
	if err := config.Parse(); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	log.Print("set upping logger")
	app_log.Setup(config.Cfg().Logger.Level)

	log.Print("connecting to minio")
	if err := app_minio.Connect(); err != nil {
		log.Fatalf("failed to connect to minio: %v", err)
	}

	log.Print("running discord bot")
	if err := discord_bot.Run(); err != nil {
		log.Fatalf("failed to run discord bot: %v", err)
	}
}
