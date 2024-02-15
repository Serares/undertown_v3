package utils

import (
	"context"
	"mime/multipart"
	"os"

	"github.com/Serares/undertown_v3/utils/env"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func UploadFilesToS3(ctx context.Context, image multipart.FileHeader, s3Client *s3.Client, imageName string) error {
	bucketName := os.Getenv(env.RAW_IMAGES_BUCKET)
	readerFile, err := image.Open()
	if err != nil {
		return err
	}
	_, err = s3Client.PutObject(
		ctx,
		&s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(imageName),
			Body:   readerFile,
		},
	)
	return err
}
