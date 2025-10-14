package repository

import (
	"Backend-RIP/internal/app/dsn"
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db                   *gorm.DB
	Interval             *IntervalRepository
	Composition_interval *CompositionIntervalRepository
	User                 *UserRepository
}

func NewRepository() (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	minioClient, err := InitMinIOClient()
	if err != nil {
		return nil, err
	}

	return &Repository{
		db:                   db,
		Interval:             NewIntervalRepository(db, minioClient),
		Composition_interval: NewCompositionIntervalRepository(db),
		User:                 NewUserRepository(db),
	}, nil
}

func CloseDBConn(r *Repository) {
	dbInstance, _ := r.db.DB()
	_ = dbInstance.Close()
}

func InitMinIOClient() (*minio.Client, error) {
	endpoint := "localhost:9000" // nginx proxy

	// Используем credentials из docker-compose
	accessKeyID := "minio"
	secretAccessKey := "minio124"
	useSSL := false

	logrus.Printf("MinIO Config - Endpoint: %s, User: %s", endpoint, accessKeyID)

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logrus.Errorf("Failed to create MinIO client: %v", err)
		return nil, fmt.Errorf("failed to create MinIO client: %v", err)
	}

	ctx := context.Background()

	// Проверим подключение
	_, err = minioClient.ListBuckets(ctx)
	if err != nil {
		logrus.Errorf("MinIO connection test failed: %v", err)
		return nil, fmt.Errorf("minio connection test failed: %v", err)
	}

	exists, err := minioClient.BucketExists(ctx, intervalImagesBucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %v", err)
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, intervalImagesBucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %v", err)
		}
		logrus.Printf("Bucket '%s' created successfully\n", intervalImagesBucket)
	}

	logrus.Printf("MinIO client initialized successfully")
	return minioClient, nil
}
