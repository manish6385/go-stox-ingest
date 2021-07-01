package main

import (
	"fmt"
	"os"
	"time"

	"github.com/g33kzone/go-stox-ingest/config"
	"github.com/g33kzone/go-stox-ingest/db"
	"github.com/g33kzone/go-stox-ingest/service"
	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Initialise logrus configuration
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(false)
	log.SetLevel(log.TraceLevel)
	log.SetOutput(os.Stdout)
}

func main() {
	var date time.Time

	// Initialize Stox app config
	conf := config.Init()

	// Initialize DB connection
	db, err := db.InitDB(&conf.PostgresConf)
	if err != nil {
		sentry.CaptureException(err)
		log.Fatal(err)
	}

	// DB Table creation - GORM Automigrate
	if err = db.CreateDBTable(); err != nil {
		sentry.CaptureException(err)
		log.Fatal(err)
	}

	bse := service.BSE{DB: db, BSEBhavCopyUrl: conf.BseBhavCopyURL}

	orders, err := bse.DB.FetchData()
	if err != nil {
		sentry.CaptureException(err)
		log.Fatal(err)
	}

	fmt.Println(orders)

	// Start BSE Bhav Copy File download
	if err = bse.InitiateDownload(date); err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
}
