package discord_bot

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"discord-recorder/internal/app_log"
	"discord-recorder/internal/app_minio"
	"discord-recorder/internal/config"
	voiceHandler "discord-recorder/internal/handler/discord/voice"
	fileService "discord-recorder/internal/service/file"
	minioService "discord-recorder/internal/service/minio"
	voiceService "discord-recorder/internal/service/voice"
	"github.com/bwmarrin/discordgo"
)

func Run() error {
	const op = "app.discord_bot.Run"

	log := app_log.Logger().With(slog.String("op", op))

	log.Info("configuring discord bot")
	discord, err := discordgo.New("Bot " + config.Cfg().Discord.Token)
	if err != nil {
		log.Error("failed to create discord bot", slog.String("error", err.Error()))
		return err
	}

	discord.StateEnabled = true
	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildVoiceStates) // todo

	setupEventListeners(discord)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Info("opening discord bot")
	if err = discord.Open(); err != nil {
		log.Error("failed to open discord bot", slog.String("error", err.Error()))
		return err
	}
	defer func(discord *discordgo.Session) {
		log.Info("closing discord bot")
		err = discord.Close()
		if err != nil {
			log.Error("failed to close discord bot", slog.String("error", err.Error()))
		}
	}(discord)

	log.Info("waiting events", slog.String("BotID", discord.State.User.ID))

	<-done

	return nil
}

func setupEventListeners(discord *discordgo.Session) {
	minioS := minioService.New(app_minio.Client(), config.Cfg().Minio.Bucket)
	fileS := fileService.New(minioS)
	voiceS := voiceService.New(config.Cfg().Discord.RecordLimit, fileS)
	voiceH := voiceHandler.New(voiceS)

	discord.AddHandler(voiceH.VoiceChannelUpdate)
}
