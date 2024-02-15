package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log/slog"
	"mime/multipart"
	"os"

	"github.com/Serares/undertown_v3/utils"
	"github.com/Serares/undertown_v3/utils/env"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/chai2010/webp"
)

var LocalAssetsRelativePath = "../../ssr/assets/uploads"

type ProcessImagesService struct {
	Log      *slog.Logger
	S3Client s3.Client
}

func New(log *slog.Logger, client s3.Client) ProcessImagesService {
	return ProcessImagesService{
		Log:      log.WithGroup("Process Images Service"),
		S3Client: client,
	}
}

func (ss *ProcessImagesService) encodeToWebP(file io.Reader) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	img, _, err := image.Decode(file)
	if err != nil {
		ss.Log.Error("error deconding the file",
			"error", err,
		)
		return nil, fmt.Errorf("error trying to decode the image")
	}

	err = webp.Encode(&buf, img, &webp.Options{Lossless: true})
	if err != nil {
		ss.Log.Error("error encoding the file",
			"error", err,
		)
	}

	return &buf, err
}

// Uploads the files to S3
// s3RawImagekey is composed of the humanReadableId of the property /humanreadableId/ImageName
func (ss *ProcessImagesService) ProcessImagesS3(ctx context.Context, s3RawImageKey, s3RawImagesBucketName string) error {
	processedImagesBucketName := os.Getenv(env.PROCESSED_IMAGES_BUCKET)
	ss.Log.Info("received image:",
		"s3key", s3RawImageKey,
	)
	s3getObject, err := ss.S3Client.GetObject(
		ctx,
		&s3.GetObjectInput{
			Bucket: &s3RawImagesBucketName,
			Key:    &s3RawImageKey,
		},
	)
	if err != nil {
		ss.Log.Error(
			"error trying to get the s3 object",
			"error", err,
		)
		return err
	}
	defer s3getObject.Body.Close()

	buf, err := ss.encodeToWebP(s3getObject.Body)
	if err != nil {
		ss.Log.Error("error trying to encode to webp",
			"error", err,
		)
	}
	processedImageKey := utils.AppendFileExtension("images/"+s3RawImageKey, "webp")
	ss.Log.Info("the processed image key", "key", processedImageKey)
	_, err = ss.S3Client.PutObject(
		ctx,
		&s3.PutObjectInput{
			Bucket: &processedImagesBucketName,
			Key:    &processedImageKey,
			Body:   buf,
		},
	)
	if err != nil {
		ss.Log.Error(
			"error uploading the file to s3",
			"error", err,
		)
	}

	// delete the image from raw images bucket if it's been successfully processed
	_, err = ss.S3Client.DeleteObject(
		ctx,
		&s3.DeleteObjectInput{
			Bucket: &s3RawImagesBucketName,
			Key:    &s3RawImageKey,
		},
	)

	return err
}

// DEPRECATED
func (ss *ProcessImagesService) ProcessImagesLocal(ctx context.Context, files []*multipart.FileHeader) ([]string, error) {
	var imageNames []string
	err := os.Mkdir(LocalAssetsRelativePath, 0755)
	if err != nil {
		ss.Log.Error("Error creating the temp dir", "err", err)
	}

	// read uploaded images
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			ss.Log.Error("error trying to read the file header")
		}
		defer file.Close()
		fileName := utils.ReplaceWhiteSpaceWithUnderscore(fileHeader.Filename)
		ss.Log.Info("the file to be encoded:", "filename", fileName)
		buf, err := ss.encodeToWebP(file)
		if err != nil {
			ss.Log.Error("error trying to encode to webp",
				"error", err,
			)
		}
		webpFileName := LocalAssetsRelativePath + "/" + fileName + ".webp"
		webpFile, err := os.Create(webpFileName)
		if err != nil {
			ss.Log.Error("error creating webFile the file",
				"error", err,
			)
		}
		defer webpFile.Close()

		_, err = webpFile.Write(buf.Bytes())
		if err != nil {
			ss.Log.Error("error writing webFile the file",
				"error", err,
			)
		}
		imageNames = append(imageNames, webpFileName)
		// remove files after beeing uploaded to s3 to not overflow to storage
	}
	return imageNames, nil
}

// DEPRECATED
func (ss *ProcessImagesService) DeleteImagesLocal(imagesNames []string) error {
	for _, imageName := range imagesNames {
		imageRelativePath := fmt.Sprintf("%s/%s", LocalAssetsRelativePath, imageName)
		if err := os.Remove(imageRelativePath); err != nil {
			ss.Log.Error("error trying to remove the file", "filePath", imageRelativePath, "error", err)
			// return fmt.Errorf("error trying to remove image %v", err) ❗this is not really needed to return because the error will generally be when files are non existent
		}
	}
	return nil
}
