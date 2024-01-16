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

type IUsersRepository interface {
	Add(ctx context.Context, userParameters lite.CreateUserParams) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*lite.User, error)
	GetByEmail(ctx context.Context, email string) (*lite.User, error)
	UpdateEmail(ctx context.Context, id string) error
	CloseDbConnection(ctx context.Context) error
}

type Users struct {
	db           *lite.Queries
	dbConnection *sql.DB
}

func NewUsersRepository() (*Users, error) {
	dbUrl, err := utils.CreateSqliteUrl()
	if err != nil {
		return nil, fmt.Errorf("error trying to create the db connection string %w", err)
	}
	db, err := sql.Open("libsql", dbUrl)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(30 * time.Minute)
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error reaching the database %w", err)
	}
	dbQueries := lite.New(db)
	return &Users{
		db:           dbQueries,
		dbConnection: db,
	}, err
}

// Implement all the functions from the interface on the Users struct
func (u *Users) Add(ctx context.Context, parameters lite.CreateUserParams) error {
	err := u.db.CreateUser(ctx, parameters)

	if err != nil {
		return fmt.Errorf("error trying to create user with id %v", parameters.ID)
	}

	return nil
}
func (u *Users) Delete(ctx context.Context, id string) error {
	err := u.db.DeleteUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.ErrorNotFound
		}
		return fmt.Errorf("%s -- %v", types.ErrorAccessingDatabase, err)
	}
	return nil
}

func (u *Users) Get(ctx context.Context, id string) (*lite.User, error) {
	user, err := u.db.GetUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrorNotFound
		}
		return nil, fmt.Errorf("%s -- %v", types.ErrorAccessingDatabase, err)
	}
	return &user, nil
}

func (u *Users) GetByEmail(ctx context.Context, email string) (*lite.User, error) {
	user, err := u.db.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrorNotFound
		}
		return nil, fmt.Errorf("%s -- %v", types.ErrorAccessingDatabase, err)
	}
	return &user, nil
}

func (u *Users) UpdateEmail(ctx context.Context, id string) error {
	return nil
}

func (u *Users) CloseDbConnection(ctx context.Context) error {
	err := u.dbConnection.Close()
	if err != nil {
		return fmt.Errorf("error on closing db connection %w", err)
	}

	return nil
}
