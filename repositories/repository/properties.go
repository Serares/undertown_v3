package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Serares/undertown_v3/repositories/repository/lite"
	"github.com/Serares/undertown_v3/repositories/repository/types"
	"github.com/Serares/undertown_v3/repositories/repository/utils"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

type IPropertiesRepository interface {
	Add(ctx context.Context, propertyParameters lite.AddPropertyParams) error
	GetById(ctx context.Context, id string, humanReadableId *string) (lite.Property, error)
	List(ctx context.Context) ([]lite.Property, error)
	DeleteByHumanReadableId(ctx context.Context, humanReadableId *string) error
	CloseDbConnection(ctx context.Context) error
	UpdateProperty(ctx context.Context, humanReadableId string, params lite.UpdatePropertyFieldsParams) error
}

type Properties struct {
	db           *lite.Queries
	dbConnection *sql.DB
}

func NewPropertiesRepo() (*Properties, error) {
	dbUrl, err := utils.CreateSqliteUrl()
	if err != nil {
		return nil, fmt.Errorf("error creating the connection string for database %w", err)
	}
	db, err := sql.Open("libsql", dbUrl)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(30 * time.Minute)
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error reaching the database %w", err)
	}
	// manage the resource usage of the db
	dbQueries := lite.New(db)

	return &Properties{
		db:           dbQueries,
		dbConnection: db,
	}, nil
}

/*
-- TODO name GetFeaturedProperties: many
-- TODO name: UpdatePropertyMetadata: exec
-- TODO  name: UpdatePropertyEnergy: exec
-- TODO name: UpdatePropertyDestination: exec
-- TODO  name: UpdatePropertyOtherUtilities: exec
--TODO  name: UpdatePropertyFurnished: exec
-- TODO name: UpdatePropertyInterior: exec
-- TODO name: UpdatePropertyHeating: exec
*/
func (d *Properties) Add(ctx context.Context, propertyParams lite.AddPropertyParams) error {
	err := d.db.AddProperty(ctx, propertyParams)

	if err != nil {
		return fmt.Errorf("error trying to insert property from user id %v error: %w", propertyParams.UserID, err)
	}

	return nil
}

// ❔
// not sure if using pointers to create nullable parameters is a good practice
// this method is used to retreive both by UUID and humanReadableId
func (d *Properties) GetById(ctx context.Context, id *string, humanReadableId *string) (lite.Property, error) {
	var property lite.Property
	var err error
	if id != nil {
		property, err = d.db.GetProperty(ctx, *id)
		if err != nil {
			return lite.Property{}, fmt.Errorf("error trying to retreive by uuid: %v, err: %v", *id, err)
		}
	}

	if humanReadableId != nil {
		property, err = d.db.GetByHumanReadableId(ctx, *humanReadableId)
		if err != nil {
			return lite.Property{}, fmt.Errorf("error trying to retreive by humanReadableId: %v, err: %v", *humanReadableId, err)
		}
	}
	return property, nil
}

func (d *Properties) List(ctx context.Context) ([]lite.Property, error) {
	properties, err := d.db.ListProperties(ctx)
	if err != nil {
		return make([]lite.Property, 0), fmt.Errorf("error listing properties: %v", err)
	}

	return properties, nil
}

func (d *Properties) ListByTransactionType(ctx context.Context, transactionType types.TransactionType) ([]lite.ListPropertiesByTransactionTypeRow, error) {
	properties, err := d.db.ListPropertiesByTransactionType(ctx, transactionType.String())
	if err != nil {
		return make([]lite.ListPropertiesByTransactionTypeRow, 0), fmt.Errorf("error listing properties by transaction type: %v", err)
	}

	return properties, nil
}

func (d *Properties) ListFeatured(ctx context.Context) ([]lite.ListFeaturedPropertiesRow, error) {
	properties, err := d.db.ListFeaturedProperties(ctx)
	if err != nil {
		return make([]lite.ListFeaturedPropertiesRow, 0), fmt.Errorf("error listing featured properties: %v", err)
	}
	return properties, nil
}

func (d *Properties) DeleteByHumanReadableId(ctx context.Context, humanReadableId *string) error {
	if humanReadableId != nil {
		err := d.db.DeletePropertyByHumanReadableId(ctx, *humanReadableId)
		if err != nil {
			return nil
		}
	}

	return fmt.Errorf("error you did not provide a human readable id")
}

func (d *Properties) UpdateProperty(ctx context.Context, humanReadableId string, params lite.UpdatePropertyFieldsParams) error {
	err := d.db.UpdatePropertyFields(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update the property features id: %s ; error: %w", humanReadableId, err)
	}

	return nil
}

func (d *Properties) CloseDbConnection(ctx context.Context) error {
	err := d.dbConnection.Close()
	if err != nil {
		return fmt.Errorf("error on closing db connection %w", err)
	}

	return nil
}
