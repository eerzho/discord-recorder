package minio

import (
	"context"
	"log/slog"
	"os"

	"discord-recorder/internal/app_log"
	"github.com/minio/minio-go/v7"
)

type Service struct {
	Client *minio.Client
	Bucket string
}

func New(client *minio.Client, bucket string) *Service {
	return &Service{Client: client, Bucket: bucket}
}

func (s *Service) Upload(filePath string) error {
	const op = "service.minio.Upload"
	log := app_log.Logger().With(slog.String("op", op))

	file, err := os.Open(filePath)
	if err != nil {
		log.Info("failed to open file", slog.String("error", err.Error()))
		return err
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Info("failed to close file", slog.String("error", err.Error()))
		}
	}(file)

	fileInfo, err := file.Stat()
	if err != nil {
		log.Error("failed to stat file", slog.String("error", err.Error()))
		return err
	}

	_, err = s.Client.PutObject(
		context.Background(),
		s.Bucket,
		fileInfo.Name(),
		file,
		fileInfo.Size(),
		minio.PutObjectOptions{
			ContentType: "application/octet-stream",
		},
	)
	if err != nil {
		log.Error("failed to upload file", slog.String("error", err.Error()))
		return err
	}

	if err = os.Remove(filePath); err != nil {
		log.Error("failed to remove file", slog.String("error", err.Error()))
	}

	return nil
}
