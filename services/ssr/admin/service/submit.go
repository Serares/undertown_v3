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
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqsTypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SubmitService struct {
	Log    *slog.Logger
	Client ISSRAdminClient
}

func NewSubmitService(log *slog.Logger, client ISSRAdminClient) *SubmitService {
	return &SubmitService{
		Log:    log.WithGroup("Submit Service"),
		Client: client,
	}
}

type PropertyFormField struct {
	Title string
}

func (s *SubmitService) Submit(r *http.Request, authToken string) (lite.Property, utils.PropertyFeatures, error) {
	piuQueuUrl := os.Getenv(env.PIU_QUEUE_URL)
	jwtSecret := os.Getenv(env.JWT_SECRET)
	var err error
	cfg, err := config.LoadDefaultConfig(r.Context())
	if err != nil {
		s.Log.Error(
			"error trying to load the lambda context",
			"error", err,
		)
		return lite.Property{}, utils.PropertyFeatures{}, err
	}

	sqsClient := sqs.NewFromConfig(cfg)

	claims, err := utils.ParseJwtWithClaims(authToken, jwtSecret)
	if err != nil {
		s.Log.Error("error trying to parse the auth token",
			"error", err,
		)
		return lite.Property{}, utils.PropertyFeatures{}, err
	}

	jsonString, humanReadableId, err := adminUtils.ParseMultipartToJson(r)
	if err != nil {
		s.Log.Error("error trying to parse the multipart/form",
			"error", err,
		)
		return lite.Property{}, utils.PropertyFeatures{}, err
	}

	messageAttributes := map[string]sqsTypes.MessageAttributeValue{
		constants.HUMAN_READABLE_ID_SQS_ATTRIBUTE: {
			DataType:    aws.String("String"),
			StringValue: aws.String(humanReadableId),
		},
		constants.USER_ID: {
			DataType:    aws.String("String"),
			StringValue: aws.String(claims.UserId),
		},
	}

	_, err = sqsClient.SendMessage(
		r.Context(),
		&sqs.SendMessageInput{
			QueueUrl:          &piuQueuUrl,
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
