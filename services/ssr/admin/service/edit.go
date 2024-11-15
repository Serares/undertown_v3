package service

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

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

type EditService struct {
	Log       *slog.Logger
	Client    ISSRAdminClient
	SQSClient *sqs.Client
	S3Client  *s3.Client
}

func NewEditService(
	log *slog.Logger,
	client ISSRAdminClient,
	sqsClient *sqs.Client,
	S3Client *s3.Client,
) *EditService {
	return &EditService{
		Log:       log.WithGroup("Edit Service"),
		Client:    client,
		SQSClient: sqsClient,
		S3Client:  S3Client,
	}
}

// TODO right now the Submit() method from the SubmitService is used to edit the property
// func (es *EditService) Post(body, humanReadableId, authToken string) error {
// }
func (es *EditService) Get(humanReadableId, authToken string) (lite.Property, []string, utils.PropertyFeatures, error) {
	getPropertyUrl := os.Getenv("GET_PROPERTY_URL")
	if es.SQSClient == nil {
		es.Log.Error("sqs client is not initialized")
		return lite.Property{}, nil, utils.PropertyFeatures{}, fmt.Errorf("error initializing the sqs client")
	}
	// have to add the human readable id to the url
	// ❗TODO
	// this might need some validations
	getPropertyBackendUrl, err := utils.AddParamToUrl(getPropertyUrl, constants.QUERY_PARAMETER_HUMANREADABLEID, humanReadableId)
	if err != nil {
		es.Log.Error("error trying to create the backend delete url", "error", err)
	}
	// process the images returned from the db
	// the images are a string separated by ;
	// return to the views the images paths as a slice of strings

	property, err := es.Client.GetProperty(getPropertyBackendUrl, authToken)
	if err != nil {
		return lite.Property{}, nil, utils.PropertyFeatures{}, err
	}

	images := strings.Split(property.Images, ";")
	images = utils.CreateImagePathList("/images/", images)
	// have to decode the property features into the utils.PropertyFeatures struct to be able to fill
	// up input values
	var propertyFeatures utils.PropertyFeatures

	err = json.Unmarshal([]byte(property.Features), &propertyFeatures)
	if err != nil {
		return lite.Property{}, nil, utils.PropertyFeatures{}, err
	}

	return property, images, propertyFeatures, nil
}

func (es *EditService) Post(r *http.Request, token, humanReadableId string) (lite.Property, utils.PropertyFeatures, error) {
	PUQueuUrl := os.Getenv(env.SQS_PU_QUEUE_URL)
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return lite.Property{}, utils.PropertyFeatures{}, err
	}
	jsonString, err := adminUtils.ParseMultipartFieldsToJson(r)
	if err != nil {
		return lite.Property{}, utils.PropertyFeatures{}, err
	}

	messageAttributes := map[string]sqsTypes.MessageAttributeValue{
		constants.HUMAN_READABLE_ID_SQS_ATTRIBUTE: {
			DataType:    aws.String("String"),
			StringValue: aws.String(humanReadableId),
		},
	}

	_, err = es.SQSClient.SendMessage(
		r.Context(),
		&sqs.SendMessageInput{
			QueueUrl:          &PUQueuUrl,
			MessageBody:       aws.String(string(jsonString)),
			MessageAttributes: messageAttributes,
		},
	)
	// TODO handle the case where iamges are removed
	if err != nil {
		es.Log.Error("error trying to send the request", "error", err)
		property, features, unmarshallingErrors := adminUtils.UnmarshalProperty(jsonString)
		// ❗
		// populating a lite.Property struct with the types.RequestProperty fields
		// because the views.Edit() component used for rerendering the last values
		// is accepting a lite.Property prop as parameter
		// TODO think if there are other ways around this hacky sittuation
		// ❗this is the ending where the fields need to be rerendered
		if unmarshallingErrors != nil {
			es.Log.Error("error trying to unmarshal the data for error response", "error", unmarshallingErrors)
			// This case the property and features will be empty struct values
			return property, features, err
		}
		return property, features, err
	}

	// this is the success ending
	return lite.Property{}, utils.PropertyFeatures{}, nil
}
