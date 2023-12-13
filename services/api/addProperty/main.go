package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Serares/undertown_v3/addProperty/internal/database"
	_ "github.com/lib/pq"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
)

var db *sql.DB // Declare a global database connection variable
var stage string

func initDatabase() error {
	stage := os.Getenv("STAGE")
	var connStr string
	switch stage {
	case "local":
		connStr = "user=localuser dbname=localdb sslmode=disable"
	case "dev":
		// Replace these values with your RDS PostgreSQL connection details
		connStr = "user=yourusername dbname=yourdatabase host=your-rds-endpoint sslmode=require password=yourpassword"
	default:
		return fmt.Errorf("unknown STAGE value: %s", stage)
	}

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	fmt.Println("Connected to the database")
	return nil
}

type POSTProperty struct {
	Title       string       `json:"name"`
	Floor       int32        `json:"floor"`
	PublishedAt sql.NullTime `json:"publishedAt"` // UTC Time
	UserID      uuid.UUID    `json:"userId"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	dbQueries := database.New(db)
	dbQueries.AddProperty(context.Background())
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error trying to initialize db repo: %s", stage)
	}
	var ppr POSTProperty
	err = json.Unmarshal([]byte(request.Body), &ppr)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error unmarshaling JSON",
		}, err
	}

	err = repo.Add(database.AddPropertyParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Title:       ppr.Title,
		Floor:       ppr.Floor,
		PublishedAt: time.Now().UTC(),
		UserID:      ppr.UserID,
	})

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error inserting property",
		}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       "POST request succeeded",
		Headers:    map[string]string{"Content-Type": "application/json"},
		StatusCode: http.StatusCreated,
	}, nil
}

func main() {
	if err := initDatabase(); err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}
	lambda.Start(handler)
}
