package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type FileStorage interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (string, error)
	GetDefaultAvatarURL() string
}

type MinioStorage struct {
	client     *minio.Client
	bucketName string
	publicURL  string
}

func NewMinioStorage(endpoint, accesKey, secretKey, bucketName, publicURL string, useSSL bool) (*MinioStorage, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accesKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return &MinioStorage{
		client:     minioClient,
		bucketName: bucketName,
		publicURL:  publicURL,
	}, nil
}

func (s *MinioStorage) UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%s/%s%s", folder, uuid.NewString(), ext)

	_, err = s.client.PutObject(ctx, s.bucketName, newFileName, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", fmt.Errorf("minio upload error: %w", err)
	}

	// http://localhost:9000/bucket-name/folder/file.jpg
	url := fmt.Sprintf("%s/%s/%s", s.publicURL, s.bucketName, newFileName)
	return url, nil
}

func (s *MinioStorage) GetDefaultAvatarURL() string {
	return fmt.Sprintf("%s/%s/avatars/default_avatar.png", s.publicURL, s.bucketName)
}
