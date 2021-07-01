package config

import (
	"fmt"
	"strconv"

	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

type Conf struct {
	BseBhavCopyURL string `env:"STOX_BSE_BHAV_URL"`
	PostgresConf   PostgresConf
	Sentry_DSN     string
}

type PostgresConf struct {
	DBInstance      string `mapstructure:"STOX_DB_INSTANCE"`
	DBName          string `env:"STOX_DB_NAME"`
	DBServer        string `env:"STOX_DB_IP_ADDR"`
	DBPort          int    `env:"STOX_DB_IP_PORT"`
	DBUser          string `env:"STOX_DB_USER"`
	DBPwd           string `env:"STOX_DB_PWD"`
	DBRetryAttempts uint8  `env:"STOX_DB_INSTANCE"`
}

func Init() *Conf {
	var conf Conf

	log.Info("Initialising Stox Ingest app configuration....")

	// viper.SetConfigName("app")
	// viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Info("No config file was found")
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}

	conf.PostgresConf.DBInstance = viper.Get("STOX_DB_INSTANCE").(string)
	conf.PostgresConf.DBName = viper.Get("STOX_DB_NAME").(string)
	conf.PostgresConf.DBServer = viper.Get("STOX_DB_IP_ADDR").(string)
	conf.PostgresConf.DBPort, _ = strconv.Atoi(viper.Get("STOX_DB_IP_PORT").(string))
	conf.PostgresConf.DBUser = viper.Get("STOX_DB_USER").(string)
	conf.PostgresConf.DBPwd = viper.Get("STOX_DB_PWD").(string)
	conf.PostgresConf.DBRetryAttempts = 2
	conf.Sentry_DSN = viper.GetString("STOX_SENTRY_DSN")
	conf.BseBhavCopyURL = viper.GetString("STOX_BSE_BHAV_URL")

	// Initliase Sentry Configuration
	initSentry(conf.Sentry_DSN)

	return &conf
}

func initSentry(sentry_dsn string) {
	if err := sentry.Init(sentry.ClientOptions{AttachStacktrace: true, Dsn: sentry_dsn}); err != nil {
		log.Fatal("Error Initializing sentry: ", "error", err.Error())
	}
}
