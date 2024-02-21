package env

// ❗TODO try renaming the variables to something more specific
// like prefixing with service name?
// Constants for env variables key names
const (
	PROCESSED_IMAGES_BUCKET               = "PROCESSED_IMAGES_BUCKET"
	RAW_IMAGES_BUCKET                     = "RAW_IMAGES_BUCKET"
	GET_PROPERTIES_URL                    = "GET_PROPERTIES_URL"
	GET_PROPERTY_URL                      = "GET_PROPERTY_URL"
	LOGIN_URL                             = "LOGIN_URL"
	DELETE_PROPERTY_URL                   = "DELETE_PROPERTY_URL"
	JWT_SECRET                            = "JWT_SECRET"
	DB_HOST                               = "DB_HOST"
	DB_NAME                               = "DB_NAME"
	DB_PROTOCOL                           = "DB_PROTOCOL"
	TURSO_DB_TOKEN                        = "TURSO_DB_TOKEN"
	SQS_PU_QUEUE_URL                      = "SQS_PU_QUEUE_URL"
	SQS_DELETE_PROCESSED_IMAGES_QUEUE_URL = "SQS_DELETE_PROCESSED_IMAGES_QUEUE_URL"
	SQS_PROCESS_RAW_IMAGES_QUEUE_URL      = "SQS_PROCESS_RAW_IMAGES_QUEUE_URL"
)
