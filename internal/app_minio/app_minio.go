package app_minio

import (
	"discord-recorder/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var client *minio.Client

func Connect() error {
	var err error

	client, err = minio.New(
		config.Cfg().Minio.Endpoint,
		&minio.Options{
			Creds: credentials.NewStaticV4(
				config.Cfg().Minio.AccessKey,
				config.Cfg().Minio.SecretKey,
				"",
			),
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func Client() *minio.Client {
	return client
}
