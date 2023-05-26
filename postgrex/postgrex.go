package postgrex

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	TimeZone string
	SSL      string
	Logger   logger.Interface
	Options
}

type Options struct {
	PreferSimpleProtocol bool
	WithoutReturning     bool
}

type poolOption func(sqlDB *sql.DB)

func New(host, port, user, password, dbname string) *Config {
	return &Config{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbname,
		Logger:   logger.Default,
		Options: Options{
			PreferSimpleProtocol: true,
		},
	}
}

func (c *Config) SetTimeZone(timezone string) *Config {
	c.TimeZone = timezone
	return c
}

func (c *Config) SetSSL(mode string) *Config {
	c.SSL = mode
	return c
}

func (c *Config) SetDB(dbname string) *Config {
	c.DBName = dbname
	return c
}

func (c *Config) SetLogger(l logger.Interface) *Config {
	c.Logger = l
	return c
}

func (c *Config) SetOptions(opts Options) *Config {
	c.Options = opts
	return c
}

func (c *Config) Connect(poolOptions ...poolOption) (*gorm.DB, error) {

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", c.Host, c.Port, c.User, c.Password, c.DBName)

	if c.TimeZone != "" {
		dsn += fmt.Sprintf(" TimeZone=%s", c.TimeZone)
	}

	if c.SSL != "" {
		dsn += fmt.Sprintf(" sslmode=%s", c.SSL)
	} else {
		dsn += fmt.Sprintf(" sslmode=%s", "disable")
	}

	setting := postgres.Config{
		DSN: dsn,
	}

	// disables implicit prepared statement usage
	if c.Options.PreferSimpleProtocol == true {
		setting.PreferSimpleProtocol = true
	}
	if c.Options.WithoutReturning == true {
		setting.WithoutReturning = true
	}

	db, _ := gorm.Open(postgres.New(setting), &gorm.Config{Logger: c.Logger})

	sqlDB, err := db.DB()
	if err != nil {
		log.Println(err)
	}

	// Options (Pool)
	for _, option := range poolOptions {
		option(sqlDB)
	}

	return db, err

}

func Pool(maxIdle, maxOpen, maxLifetime int) poolOption {
	return func(sqlDB *sql.DB) {
		sqlDB.SetMaxIdleConns(maxIdle)
		sqlDB.SetMaxOpenConns(maxOpen)
		sqlDB.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Second)
	}
}

func Level(lv string) logger.LogLevel {
	switch lv {
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Silent
	}
}
