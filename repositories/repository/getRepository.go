package repostiroy

import (
	"context"
	"fmt"
	"os"

	"github.com/Serares/undertown_v3/repositories/repository/postgres"
	"github.com/Serares/undertown_v3/repositories/repository/psql"
	"github.com/Serares/undertown_v3/repositories/repository/sqlite"
	"github.com/Serares/undertown_v3/repositories/repository/types"
	"github.com/Serares/undertown_v3/repositories/repository/utils"
	"github.com/google/uuid"
)

var DRIVER = os.Getenv("DB_DRIVER")

type IUsersRepository interface {
	Add(ctx context.Context, userParameters interface{}) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*psql.User, error)
	GetByEmail(ctx context.Context, email string) (*psql.User, error)
	UpdateEmail(ctx context.Context, id string) error
	CloseDbConnection(ctx context.Context) error
	// Used for tests
	// tryed it and it's not doing the thing
	// maybe try again later
	// BeginTx(ctx context.Context) (*sql.Tx, error)
}

type IPropertiesRepository interface {
	Add(ctx context.Context, propertyParameters psql.AddPropertyParams) error
	GetById(ctx context.Context, id *uuid.NullUUID, humanReadableId *string) (psql.Property, error)
	List(ctx context.Context) ([]psql.Property, error)
	DeleteByHumanReadableId(ctx context.Context, humanReadableId *string) error
	CloseDbConnection(ctx context.Context) error
}

func GetUserRepository(ctx context.Context) (IUsersRepository, error) {
	if DRIVER == types.PSQL_DRIVER {
		dbUrl, err := utils.CreatePsqlUrl(ctx)
		if err != nil {
			return nil, fmt.Errorf("error trying to initialize postgres repository %w", err)
		}
		ur, err := postgres.NewUsersRepository(dbUrl)
		if err != nil {
			return nil, err
		}
		return ur, err
	} else {
		dbUrl, err := utils.GetSqliteUrl()
		if err != nil {
			return nil, fmt.Errorf("error creating the sqlite url %w", err)
		}
		if err != nil {
			ur, err := sqlite.NewUsersRepository(dbUrl)
			return nil, fmt.Errorf("error initialting the sqlite users repository %v", err)
		}

		return ur, nil
	}
}

func GetPropertiesRepository() (IPropertiesRepository, error) {
	return nil, fmt.Errorf("not yet implemented")
}
