package common

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Ved1234"
	dbname   = "restaurant"
)

var DB *sqlx.DB

func DbConnection() {
	postgresql := fmt.Sprintf("host=%s port=%d user= %s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sqlx.Open("postgres", postgresql)
	if err != nil {
		fmt.Println("err")
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	DB = db
	fmt.Println("database connection is done")
	driver, driverErr := postgres.WithInstance(db.DB, &postgres.Config{})
	if driverErr != nil {

	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres://postgres:Ved1234@localhost:5432/restaurant?sslmode=disable&search_path=public",
		driver)
	if err != nil {
		log.Fatal(err)
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("error in migration ", err)
	}
}

// Tx provides the transaction wrapper
func Tx(fn func(tx *sqlx.Tx) error) error {
	tx, err := DB.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start a transaction: %+v", err)
	}
	defer func() {
		if err != nil {
			if rollBackErr := tx.Rollback(); rollBackErr != nil {
				logrus.Errorf("failed to rollback tx: %s", rollBackErr)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			logrus.Errorf("failed to commit tx: %s", commitErr)
		}
	}()
	err = fn(tx)
	return err
}
