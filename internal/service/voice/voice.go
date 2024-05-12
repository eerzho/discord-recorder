package voice

import (
	"log/slog"
	"sync"
	"time"

	"discord-recorder/internal/app_log"
	"github.com/bwmarrin/discordgo"
)

type FileService interface {
	StartRecording(c chan *discordgo.Packet, disconnectChan chan struct{}, channelID string)
}

type connectionState struct {
	vc           *discordgo.VoiceConnection
	disconnectCh chan struct{}
	userCount    int
}

type Service struct {
	mu          sync.Mutex
	limitSec    time.Duration
	fileService FileService
	connections map[string]*connectionState
}

func New(limitSec int, fileService FileService) *Service {
	return &Service{
		limitSec:    time.Duration(limitSec),
		fileService: fileService,
		connections: make(map[string]*connectionState),
	}
}

func (s *Service) Connect(session *discordgo.Session, guildID, channelID string, mute, deaf bool) {
	const op = "service.voice.Connect"
	log := app_log.Logger().With(slog.String("op", op))

	conn := s.get(channelID)
	if conn != nil {
		conn.userCount++
		log.Info("already connected to voice channel")
		return
	}

	vc, err := session.ChannelVoiceJoin(guildID, channelID, mute, deaf)
	if err != nil {
		log.Error("failed to connect to voice channel", slog.String("error", err.Error()))
		return
	}

	conn = &connectionState{vc: vc, disconnectCh: make(chan struct{}), userCount: 1}
	s.add(channelID, conn)

	log.Info("connected and starting recording")
	go s.fileService.StartRecording(vc.OpusRecv, conn.disconnectCh, channelID)

	t := time.NewTimer(s.limitSec * time.Second)
	select {
	case <-conn.disconnectCh:
		log.Info("disconnect signal received")
	case <-t.C:
		log.Info("timer finished, disconnecting")
		close(conn.disconnectCh)
	}

	s.cleanup(channelID)
}

func (s *Service) Disconnect(channelID string) {
	const op = "service.voice.Disconnect"
	log := app_log.Logger().With(slog.String("op", op))

	conn := s.get(channelID)
	if conn == nil {
		log.Info("no active connection in this channel")
		return
	}

	if conn.userCount > 1 {
		conn.userCount--
		log.Info("channel is active", slog.Int("userCount", conn.userCount))
		return
	}

	log.Info("disconnecting from voice channel")
	close(conn.disconnectCh)
	if err := conn.vc.Disconnect(); err != nil {
		log.Error("failed to disconnect from voice channel", slog.String("error", err.Error()))
	}
	s.cleanup(channelID)
}

func (s *Service) cleanup(channelID string) {
	const op = "service.voice.cleanup"
	log := app_log.Logger().With(slog.String("op", op))

	conn := s.get(channelID)

	if conn != nil {
		if conn.vc != nil {
			if err := conn.vc.Disconnect(); err != nil {
				log.Error("failed to disconnect from voice channel", slog.String("error", err.Error()))
			}
			conn.vc.Close()
		}
		s.remove(channelID)
	}
	log.Info("cleaned up connection and voice channel")
}

func (s *Service) add(channelID string, conn *connectionState) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.connections[channelID] = conn
}

func (s *Service) remove(channelID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.connections, channelID)
}

func (s *Service) get(channelID string) *connectionState {
	s.mu.Lock()
	defer s.mu.Unlock()

	conn, ok := s.connections[channelID]
	if !ok {
		return nil
	}

	return conn
}
