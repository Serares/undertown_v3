package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Serares/undertown_v3/repositories/repository/psql"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Properties struct {
	db           *psql.Queries
	dbConnection *sql.DB
}

func NewPropertiesRepo(dbUrl string) (*Properties, error) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}
	// manage the resource usage of the db
	dbQueries := psql.New(db)

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
func (d *Properties) Add(ctx context.Context, propertyParams psql.AddPropertyParams) error {
	err := d.db.AddProperty(ctx, propertyParams)

	if err != nil {
		return fmt.Errorf("error trying to insert property with id %v from user id %v", propertyParams.ID, propertyParams.UserID)
	}

	return nil
}

// ❔
// not sure if using pointers to create nullable parameters is a good practice
// this method is used to retreive both by UUID and humanReadableId
func (d *Properties) GetById(ctx context.Context, id *uuid.UUID, humanReadableId *string) (psql.Property, error) {
	var property psql.Property
	var err error
	if id != nil {
		property, err = d.db.GetProperty(ctx, *id)
		if err != nil {
			return psql.Property{}, fmt.Errorf("error trying to retreive by uuid: %v, err: %v", *id, err)
		}
	}

	if humanReadableId != nil {
		property, err = d.db.GetByHumanReadableId(ctx, *humanReadableId)
		if err != nil {
			return psql.Property{}, fmt.Errorf("error trying to retreive by humanReadableId: %v, err: %v", *humanReadableId, err)
		}
	}
	return property, nil
}

func (d *Properties) List(ctx context.Context) ([]psql.Property, error) {
	properties, err := d.db.ListProperties(ctx)
	if err != nil {
		return make([]psql.Property, 0), fmt.Errorf("error listing properties: %v", err)
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
