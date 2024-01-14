package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Serares/undertown_v3/repositories/repository/lite"
	_ "github.com/lib/pq"
)

type IPropertiesRepository interface {
	Add(ctx context.Context, propertyParameters lite.AddPropertyParams) error
	GetById(ctx context.Context, id string, humanReadableId *string) (lite.Property, error)
	List(ctx context.Context) ([]lite.Property, error)
	DeleteByHumanReadableId(ctx context.Context, humanReadableId *string) error
	CloseDbConnection(ctx context.Context) error
}

type Properties struct {
	db           *lite.Queries
	dbConnection *sql.DB
}

func NewPropertiesRepo(dbUrl string) (*Properties, error) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
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
		return fmt.Errorf("error trying to insert property with id %v from user id %v", propertyParams.ID, propertyParams.UserID)
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

func (d *Properties) DeleteByHumanReadableId(ctx context.Context, humanReadableId *string) error {
	if humanReadableId != nil {
		err := d.db.DeletePropertyByHumanReadableId(ctx, *humanReadableId)
		if err != nil {
			return nil
		}
	}

	return fmt.Errorf("error you did not provide a human readable id")
}

func (d *Properties) CloseDbConnection(ctx context.Context) error {
	err := d.dbConnection.Close()
	if err != nil {
		return fmt.Errorf("error on closing db connection %w", err)
	}

	return nil
}
