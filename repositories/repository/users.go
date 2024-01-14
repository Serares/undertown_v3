package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Serares/undertown_v3/repositories/repository/lite"
	_ "github.com/lib/pq"
)

type IUsersRepository interface {
	Add(ctx context.Context, userParameters lite.CreateUserParams) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*lite.User, error)
	GetByEmail(ctx context.Context, email string) (*lite.User, error)
	UpdateEmail(ctx context.Context, id string) error
	CloseDbConnection(ctx context.Context) error
	// Used for tests
	// tryed it and it's not doing the thing
	// maybe try again later
	// BeginTx(ctx context.Context) (*sql.Tx, error)
}

type Users struct {
	db           *lite.Queries
	dbConnection *sql.DB
}

func NewUsersRepository(dbUrl string) (*Users, error) {
	db, err := sql.Open("libsql", dbUrl)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
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
			return ErrorNotFound
		}
		return fmt.Errorf("%s -- %v", ErrorAccessingDatabase, err)
	}
	return nil
}

func (u *Users) Get(ctx context.Context, id string) (*lite.User, error) {
	user, err := u.db.GetUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrorNotFound
		}
		return nil, fmt.Errorf("%s -- %v", ErrorAccessingDatabase, err)
	}
	return &user, nil
}

func (u *Users) GetByEmail(ctx context.Context, email string) (*lite.User, error) {
	user, err := u.db.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrorNotFound
		}
		return nil, fmt.Errorf("%s -- %v", ErrorAccessingDatabase, err)
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

// func (u *Users) BeginTx(ctx context.Context) (*sql.Tx, error) {
// 	return u.dbConnection.Begin()
// You'll have to use the tx struct and run tx.Rollback()
// }
