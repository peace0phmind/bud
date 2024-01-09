package infra

import (
	"embed"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/peace0phmind/bud/factory"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"strings"
	"time"
)

var _db = factory.Singleton[gorm.DB]().SetInitOnce(initDB).MustBuilder()

func initDB() (db *gorm.DB, err error) {
	var logLevel logger.LogLevel

	conf := _config()

	strLogLevel := strings.ToLower(conf.DBLogLevel)

	switch strLogLevel {
	case "info":
		logLevel = logger.Info
	case "warn":
		logLevel = logger.Warn
	case "error":
		logLevel = logger.Error
	case "silent":
		logLevel = logger.Silent
	default:
		logLevel = logger.Error
	}

	newLogger := logger.New(
		log.New(logrus.StandardLogger().Out, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			LogLevel:                  logLevel, // Log level
			IgnoreRecordNotFoundError: false,    // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,    // Disable color
		},
	)

	db, err = gorm.Open(mysql.Open(_config().DSN()), &gorm.Config{
		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		logrus.Fatalf("initOnce db fail: %v", err)
		return nil, err
	}
	db.Logger.LogMode(logger.Info)

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	// fix: closing bad idle connection: EOF
	sqlDB.SetConnMaxIdleTime(1 * time.Hour)
	sqlDB.SetConnMaxLifetime(6 * time.Hour)

	logrus.Info("create db success")
	return
}

func DoMigrate(migrationsFS embed.FS) {
	d, err := iofs.New(migrationsFS, "migration")
	if err != nil {
		logrus.Fatal(err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, "mysql://"+_config().DSN())
	if err != nil {
		logrus.Fatal(err)
	}

	if err = m.Up(); err != nil {
		logrus.Warn(err)
	}
}
