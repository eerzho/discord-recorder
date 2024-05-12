package voice

import (
	"log/slog"

	"discord-recorder/internal/app_log"
	"github.com/bwmarrin/discordgo"
)

type service interface {
	Connect(session *discordgo.Session, guildID, channelID string, mute, deaf bool)
	Disconnect(channelID string)
}

type Handler struct {
	service service
}

func New(service service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) VoiceChannelUpdate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	const op = "handler.discord.voice.VoiceChannelUpdate"
	log := app_log.Logger().With(
		slog.String("method", op),
		slog.String("UserID", v.UserID),
		slog.String("GuildID", v.GuildID),
		slog.String("ChannelID", v.ChannelID),
	)

	if s.State.User.ID == v.UserID {
		return
	}

	if v.GuildID != "861893956078272513" {
		return
	}

	if v.ChannelID != "" {
		log.Info("connecting to voice channel")
		go h.service.Connect(s, v.GuildID, v.ChannelID, true, false)
	} else if v.ChannelID == "" && v.BeforeUpdate != nil && v.BeforeUpdate.ChannelID != "" {
		log.Info("disconnecting from voice channel")
		go h.service.Disconnect(v.BeforeUpdate.ChannelID)
	}
}
