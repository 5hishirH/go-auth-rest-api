package minio

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/5hishirH/go-auth-rest-api.git/internal/shared/storage/filestore"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	minioClient *minio.Client
	bucketName  string
}

var _ filestore.FileStore = (*Client)(nil)

// New initializes the connection to MinIO
func New(endpoint, accessKey, secretKey, bucketName string, useSSL bool) (*Client, error) {
	// Initialize MinIO client object
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	// Ensure bucket exists (Optional, but good for local dev)
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
		log.Printf("Bucket %s created successfully", bucketName)
	}

	return &Client{
		minioClient: minioClient,
		bucketName:  bucketName,
	}, nil
}

// UploadFile uploads a stream of data to MinIO and returns the object name
func (c *Client) Upload(ctx context.Context, objectName string, reader io.Reader, fileSize int64, contentType string) error {
	_, err := c.minioClient.PutObject(ctx, c.bucketName, objectName, reader, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})

	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

func (c *Client) Delete(ctx context.Context, objectName string) error {
	return c.minioClient.RemoveObject(ctx, c.bucketName, objectName, minio.RemoveObjectOptions{})
}
