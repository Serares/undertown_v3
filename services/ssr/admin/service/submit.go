package service

import (
	"log/slog"
	"net/http"
	"os"

	adminUtils "github.com/Serares/ssr/admin/utils"
	"github.com/Serares/undertown_v3/repositories/repository/lite"
	"github.com/Serares/undertown_v3/utils"
	"github.com/Serares/undertown_v3/utils/constants"
	"github.com/Serares/undertown_v3/utils/env"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqsTypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SubmitService struct {
	Log       *slog.Logger
	Client    ISSRAdminClient
	SQSClient *sqs.Client
	S3Client  *s3.Client
}

func NewSubmitService(
	log *slog.Logger,
	client ISSRAdminClient,
	sqsClient *sqs.Client,
	s3Client *s3.Client,
) *SubmitService {
	return &SubmitService{
		Log:       log.WithGroup("Submit Service"),
		Client:    client,
		SQSClient: sqsClient,
		S3Client:  s3Client,
	}
}

type PropertyFormField struct {
	Title string
}

func (s *SubmitService) Submit(r *http.Request, authToken string) (lite.Property, utils.PropertyFeatures, error) {
	PUQueuUrl := os.Getenv(env.SQS_PU_QUEUE_URL)
	jwtSecret := os.Getenv(env.JWT_SECRET)
	var err error

	if s.SQSClient == nil {
		s.Log.Error("sqs client is not initialized")
		return lite.Property{}, utils.PropertyFeatures{}, err
	}

	if s.S3Client == nil {
		s.Log.Error("s3 client is not initialized")
		return lite.Property{}, utils.PropertyFeatures{}, err
	}

	claims, err := utils.ParseJwtWithClaims(authToken, jwtSecret)
	if err != nil {
		s.Log.Error("error trying to parse the auth token",
			"error", err,
		)
		return lite.Property{}, utils.PropertyFeatures{}, err
	}

	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		return lite.Property{}, utils.PropertyFeatures{}, err
	}

	// ⚠️😒
	// This pattern is a massive 💩 should move to json
	jsonString, err := adminUtils.ParseMultipartFieldsToJson(r)
	if err != nil {
		s.Log.Error("error trying to parse the multipart/form",
			"error", err,
		)
		return lite.Property{}, utils.PropertyFeatures{}, err
	}

	messageAttributes := map[string]sqsTypes.MessageAttributeValue{
		constants.USER_ID: {
			DataType:    aws.String("String"),
			StringValue: aws.String(claims.UserId),
		},
	}

	_, err = s.SQSClient.SendMessage(
		r.Context(),
		&sqs.SendMessageInput{
			QueueUrl:          &PUQueuUrl,
			MessageBody:       aws.String(string(jsonString)),
			MessageAttributes: messageAttributes,
		},
	)
	if err != nil {
		s.Log.Error("error trying to send the sqs message", "error", err)
		property, features, unmarshallingErrors := adminUtils.UnmarshalProperty(jsonString)
		// ❗
		// populating a lite.Property struct with the types.RequestProperty fields
		// because the views.Edit() component used for rerendering the last values
		// is accepting a lite.Property prop as parameter
		// TODO think if there are other ways around this hacky sittuation
		// ❗this is the ending where the fields need to be rerendered
		if unmarshallingErrors != nil {
			s.Log.Error("error trying to unmarshal the data for error response", "error", unmarshallingErrors)
			// This case the property and features will be empty struct values
			return property, features, err
		}
		return property, features, err
	}

	// this is the success ending
	return lite.Property{}, utils.PropertyFeatures{}, nil
}
