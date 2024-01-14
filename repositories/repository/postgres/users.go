package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Serares/undertown_v3/repositories/repository/psql"
	"github.com/Serares/undertown_v3/repositories/repository/types"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Users struct {
	db           *psql.Queries
	dbConnection *sql.DB
}

func NewUsersRepository(dbUrl string) (*Users, error) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	dbQueries := psql.New(db)
	return &Users{
		db:           dbQueries,
		dbConnection: db,
	}, err
}

// Implement all the functions from the interface on the Users struct
func (u *Users) Add(ctx context.Context, parameters psql.CreateUserParams) error {
	err := u.db.CreateUser(ctx, parameters)

	if err != nil {
		return fmt.Errorf("error trying to create user with id %v", parameters.ID)
	}

	return nil
}
func (u *Users) Delete(ctx context.Context, id uuid.UUID) error {
	err := u.db.DeleteUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.ErrorNotFound
		}
		return fmt.Errorf("%s -- %v", types.ErrorAccessingDatabase, err)
	}
	return nil
}

func (u *Users) Get(ctx context.Context, id uuid.UUID) (*psql.User, error) {
	user, err := u.db.GetUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrorNotFound
		}
		return nil, fmt.Errorf("%s -- %v", types.ErrorAccessingDatabase, err)
	}
	return &user, nil
}

func (u *Users) GetByEmail(ctx context.Context, email string) (*psql.User, error) {
	user, err := u.db.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrorNotFound
		}
		return nil, fmt.Errorf("%s -- %v", types.ErrorAccessingDatabase, err)
	}
	return &user, nil
}

func (u *Users) UpdateEmail(ctx context.Context, id uuid.UUID) error {
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
