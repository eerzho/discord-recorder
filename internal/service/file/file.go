package file

import (
	"fmt"
	"log/slog"
	"time"

	"discord-recorder/internal/app_log"
	"github.com/bwmarrin/discordgo"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/pion/webrtc/v3/pkg/media/oggwriter"
)

type MinioService interface {
	Upload(filePath string) error
}

type Service struct {
	minioService MinioService
}

func New(minioService MinioService) *Service {
	return &Service{minioService: minioService}
}

func (s *Service) StartRecording(c chan *discordgo.Packet, disconnectChan chan struct{}, channelID string) {
	const op = "service.file.StartRecording"
	log := app_log.Logger().With(slog.String("op", op))

	oggFilePath := fmt.Sprintf("tmp/%s_%d.ogg", channelID, time.Now().Unix())
	files := make(map[uint32]media.Writer)

	defer s.closeAllFiles(files, oggFilePath)

	for {
		select {
		case p, ok := <-c:
			if !ok {
				log.Info("channel closed, stopping recording")
				return
			}
			if file, err := s.getOrCreateFile(files, oggFilePath, p.SSRC); err == nil {
				s.writePacket(file, p)
			}
		case <-disconnectChan:
			log.Info("disconnect signal received, stopping recording")
			return
		}
	}
}

func (s *Service) getOrCreateFile(files map[uint32]media.Writer, filePath string, ssrc uint32) (media.Writer, error) {
	const op = "service.file.getOrCreateFile"
	log := app_log.Logger().With(slog.String("op", op))

	if file, exists := files[ssrc]; exists {
		return file, nil
	}
	file, err := oggwriter.New(filePath, 48000, 2)
	if err != nil {
		log.Error("failed to create ogg file", slog.String("error", err.Error()))
		return nil, err
	}
	files[ssrc] = file
	return file, nil
}

func (s *Service) writePacket(file media.Writer, p *discordgo.Packet) {
	const op = "service.file.writePacket"
	log := app_log.Logger().With(slog.String("op", op))

	rtpPacket := &rtp.Packet{
		Header: rtp.Header{
			Version:        2,
			PayloadType:    0x78,
			SequenceNumber: p.Sequence,
			Timestamp:      p.Timestamp,
			SSRC:           p.SSRC,
		},
		Payload: p.Opus,
	}
	if err := file.WriteRTP(rtpPacket); err != nil {
		log.Error("failed to write to ogg file", slog.String("error", err.Error()))
	}
}

func (s *Service) closeAllFiles(files map[uint32]media.Writer, filePath string) {
	const op = "service.file.closeAllFiles"
	log := app_log.Logger().With(slog.String("op", op))

	for _, f := range files {
		if err := f.Close(); err != nil {
			log.Error("failed to close file", slog.String("error", err.Error()))
		}
		if err := s.minioService.Upload(filePath); err != nil {
			log.Error("failed to upload file", slog.String("error", err.Error()))
		}
	}
}
