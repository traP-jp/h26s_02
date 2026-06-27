package s3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/traP-jp/h26s_02/storage"
)

const defaultPresignExpiration = 15 * time.Minute

var ErrImageNotFound = errors.New("image not found")

var _ storage.Image = (*Image)(nil)

type Image struct {
	client        *awss3.Client
	presignClient *awss3.PresignClient
	bucketName    string
	presignExpire time.Duration
}

func NewImage() (*Image, error) {
	bucketName := os.Getenv("S3_BUCKET_NAME")
	if bucketName == "" {
		return nil, errors.New("S3_BUCKET_NAME is required")
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	loadOptions := []func(*config.LoadOptions) error{
		config.WithRegion(region),
	}

	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if accessKeyID != "" && secretAccessKey != "" {
		loadOptions = append(loadOptions, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, ""),
		))
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), loadOptions...)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	endpoint := os.Getenv("S3_ENDPOINT")
	client := awss3.NewFromConfig(cfg, func(o *awss3.Options) {
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
		o.UsePathStyle = true
	})

	presignClient := awss3.NewPresignClient(client)
	if endpoint == "http://s3:9000" {
		localClient := awss3.NewFromConfig(cfg, func(o *awss3.Options) {
			if endpoint != "" {
				o.BaseEndpoint = aws.String("http://localhost:9000")
			}
			o.UsePathStyle = true
		})
		presignClient = awss3.NewPresignClient(localClient)
	}

	return &Image{
		client:        client,
		presignClient: presignClient,
		bucketName:    bucketName,
		presignExpire: defaultPresignExpiration,
	}, nil
}

func (i *Image) nameToKey(name string) *string {
	return aws.String(fmt.Sprintf("image/%s", name))
}

func (i *Image) SaveImage(ctx context.Context, name string, image io.Reader) error {
	if name == "" {
		return errors.New("image name is required")
	}
	if image == nil {
		return errors.New("image reader is required")
	}

	_, err := i.client.PutObject(ctx, &awss3.PutObjectInput{
		Bucket: aws.String(i.bucketName),
		Key:    i.nameToKey(name),
		Body:   image,
	})
	if err != nil {
		return fmt.Errorf("put image to s3: %w", err)
	}

	return nil
}

func (i *Image) GetTemporalyUrl(ctx context.Context, name string) (string, error) {
	if name == "" {
		return "", errors.New("image name is required")
	}

	key := i.nameToKey(name)
	_, err := i.client.HeadObject(ctx, &awss3.HeadObjectInput{
		Bucket: aws.String(i.bucketName),
		Key:    key,
	})
	if err != nil {
		if isNotFoundError(err) {
			return "", storage.ErrImageNotFound
		}

		return "", fmt.Errorf("head image from s3: %w", err)
	}

	result, err := i.presignClient.PresignGetObject(ctx, &awss3.GetObjectInput{
		Bucket: aws.String(i.bucketName),
		Key:    key,
	}, func(options *awss3.PresignOptions) {
		options.Expires = i.presignExpire
	})
	if err != nil {
		return "", fmt.Errorf("presign get image url: %w", err)
	}

	return result.URL, nil
}

func isNotFoundError(err error) bool {
	var apiErr smithy.APIError
	if !errors.As(err, &apiErr) {
		return false
	}

	return apiErr.ErrorCode() == "NotFound" || apiErr.ErrorCode() == "NoSuchKey" || apiErr.ErrorCode() == "404"
}
