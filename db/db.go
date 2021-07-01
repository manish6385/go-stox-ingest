package db

import (
	"fmt"

	"github.com/g33kzone/go-stox-ingest/config"
	models "github.com/g33kzone/go-stox-ingest/models/entity"
	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

//Database struct
type Database struct {
	DbConn *gorm.DB
}

// Datastore is an interface to the backend datastore
type Datastore interface {
	CreateDBTable() error
	UploadBSEBhavCopy(bse_bhav []models.BSE_BHAV) error
	FetchData() ([]models.Order, error)
}

// InitializeDB creates a DB connection from the provided configuration
func InitDB(pgConf *config.PostgresConf) (Datastore, error) {

	dbDSN := fmt.Sprintf("postgres://%s@%s:%d/%s?sslmode=disable", pgConf.DBUser, pgConf.DBServer, pgConf.DBPort, pgConf.DBName)

	db, err := gorm.Open(postgres.New(postgres.Config{DSN: dbDSN}), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		sentry.CaptureException(err)
		return nil, err
	}

	log.Info("Database connection successful")

	return &Database{db}, nil
}

// Init - Create DB tables
func (db *Database) CreateDBTable() error {
	if err := db.DbConn.AutoMigrate(&models.Holiday{}); err != nil {
		return err
	}

	if err := db.DbConn.AutoMigrate(&models.BSE_BHAV{}); err != nil {
		return err
	}

	if err := db.DbConn.AutoMigrate(&models.Order{}); err != nil {
		return err
	}

	return nil
}
