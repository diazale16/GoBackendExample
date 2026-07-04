package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type Config struct {
	ProjectID string
	Endpoint  string
	Region    string
	AccessKey string
	SecretKey string
	Bucket    string
}

type S3Client struct {
	client    *s3.Client
	bucket    string
	projectID string
}

type UploadResult struct {
	Key      string
	URL      string
	MimeType string
}

func New(cfg Config) (*S3Client, error) {
	ctx := context.Background()

	var awsCfg aws.Config
	var err error

	creds := credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")

	if cfg.Endpoint != "" {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:               cfg.Endpoint,
				SigningRegion:     cfg.Region,
				HostnameImmutable: true,
			}, nil
		})

		awsCfg, err = awsconfig.LoadDefaultConfig(ctx,
			awsconfig.WithRegion(cfg.Region),
			awsconfig.WithCredentialsProvider(creds),
			awsconfig.WithEndpointResolverWithOptions(customResolver),
		)
	} else {
		awsCfg, err = awsconfig.LoadDefaultConfig(ctx,
			awsconfig.WithRegion(cfg.Region),
			awsconfig.WithCredentialsProvider(creds),
		)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg)

	log.Printf("S3 client initialized (endpoint: %s, bucket: %s)", cfg.Endpoint, cfg.Bucket)

	return &S3Client{
		client:    client,
		bucket:    cfg.Bucket,
		projectID: cfg.ProjectID,
	}, nil
}

func (s *S3Client) Upload(ctx context.Context, data io.Reader, filename string, contentType string) (*UploadResult, error) {
	key := generateKey(filename)

	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        data,
		ContentType: aws.String(contentType),
	}

	_, err := s.client.PutObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	url := s.buildURL(key)
	log.Printf("File uploaded: %s -> %s", filename, key)

	return &UploadResult{
		Key:      key,
		URL:      url,
		MimeType: contentType,
	}, nil
}

func (s *S3Client) Delete(ctx context.Context, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	_, err := s.client.DeleteObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	log.Printf("File deleted: %s", key)
	return nil
}

func (s *S3Client) GetURL(key string) string {
	return s.buildURL(key)
}

func (s *S3Client) buildURL(key string) string {
	return fmt.Sprintf("https://%s.storage.supabase.co/storage/v1/object/authenticated/%s/%s", s.projectID, s.bucket, key)
}

func generateKey(filename string) string {
	id := uuid.New().String()
	return fmt.Sprintf("uploads/%s/%s", time.Now().Format("2006/01/02"), id+"_"+filename)
}

func (s *S3Client) BuildURL(key string) string {
	return s.buildURL(key)
}