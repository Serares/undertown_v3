package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Serares/undertown_v3/internal/database"
	_ "github.com/lib/pq"
)

// type Repository interface {
// 	Create() (int64, error)
// 	Update() error
// 	ByID(id int64) (, error)
// 	Last() (, error)
// 	Breaks(n int) (, error)
// 	CategorySummary(day time.Time, filter string) (time.Duration, error)
// }

type dbRepo struct {
	db *database.Queries
}

func NewPSQLRepo(dbUrl string) (*dbRepo, error) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}
	// manage the resource usage of the db
	dbQueries := database.New(db)

	return &dbRepo{
		db: dbQueries,
	}, nil
}

func (d *dbRepo) Add(propertyParams database.AddPropertyParams) error {
	err := d.db.AddProperty(context.Background(), propertyParams)

	if err != nil {
		return fmt.Errorf("error trying to insert property with id %s from user id %s", propertyParams.ID, propertyParams.UserID)
	}

	return nil
}
