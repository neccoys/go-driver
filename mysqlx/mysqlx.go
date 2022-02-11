package mysqlx

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net/url"
	"reflect"
	"strings"
	"time"
)

type Config struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	AppendParameter string
	Logger          logger.Interface
	DefaultParameters
	Options
}

type Options struct {
	SkipInitializeWithVersion bool // 根据当前 MySQL 版本自动配置
	DefaultStringSize         uint // string 类型字段的默认长度
	DefaultDatetimePrecision  int
	DisableDatetimePrecision  bool // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
	DontSupportRenameIndex    bool // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
	DontSupportRenameColumn   bool // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
	DontSupportForShareClause bool
}

type DefaultParameters struct {
	Charset   string
	ParseTime string
	Loc       string // parseTime=true
}

type CustomParameter struct {
	Timeout           string
	ReadTimeout       string
	WriteTimeout      string
	CheckConnLiveness string // package default:true
	Collation         string
	ClientFoundRows   string
	ColumnsWithAlias  string
	Tls               string
}

func NewParamsmeter() CustomParameter {
	return CustomParameter{
		CheckConnLiveness: "true",
		Timeout:           "180s",
	}
}

func (c *Config) SetAppendParameter(params CustomParameter) *Config {
	tp := reflect.TypeOf(params)
	val := reflect.ValueOf(params)

	num := tp.NumField()
	paramsArr := make([]string, 0, num)
	for i := 0; i < num; i++ {
		key := tp.Field(i).Name

		v := fmt.Sprintf("%v", val.Field(i).Interface())
		if v != "" {
			paramsArr = append(paramsArr, fmt.Sprintf("%s=%v", strings.ToLower(key[:1])+key[1:], v))
		}

	}

	c.AppendParameter = "&" + strings.Join(paramsArr, "&")

	return c
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
		DefaultParameters: DefaultParameters{
			Charset:   "utf8mb4",
			ParseTime: "true",
			Loc:       "UTC",
		},
		Options: Options{
			SkipInitializeWithVersion: false,
			DontSupportForShareClause: true,
			DontSupportRenameColumn:   true,
			DontSupportRenameIndex:    true, // 5.7以下不支援
			DisableDatetimePrecision:  true, // 5.6以下不支援
		},
	}
}

func (c *Config) SetParseTime(p bool) *Config {
	c.DefaultParameters.ParseTime = fmt.Sprintf("%v", p)
	return c
}

func (c *Config) SetLoc(loc string) *Config {
	c.DefaultParameters.ParseTime = "true"
	c.DefaultParameters.Loc = loc
	return c
}

func (c *Config) SetCharset(charset string) *Config {
	c.DefaultParameters.Charset = charset
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

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", c.User, c.Password, c.Host, c.Port, c.DBName, c.DefaultParameters.Charset)

	if c.DefaultParameters.ParseTime != "" {
		dsn += fmt.Sprintf("&parseTime=%s", c.DefaultParameters.ParseTime)
	}
	if c.DefaultParameters.Loc != "" {
		dsn += fmt.Sprintf("&loc=%s", url.QueryEscape(c.DefaultParameters.Loc))
	}

	if c.AppendParameter != "" {
		dsn += c.AppendParameter
	}

	setting := mysql.Config{
		DSN: dsn,
	}

	if c.Options.SkipInitializeWithVersion == true {
		setting.SkipInitializeWithVersion = true

		if c.Options.DontSupportRenameIndex == true {
			setting.DontSupportRenameIndex = true
		} else {
			setting.DontSupportRenameIndex = false
		}
		if c.Options.DontSupportRenameColumn == true {
			setting.DontSupportRenameColumn = true
		} else {
			setting.DontSupportRenameColumn = false
		}
		if c.Options.DontSupportForShareClause == true {
			setting.DontSupportForShareClause = true
		} else {
			setting.DontSupportForShareClause = false
		}
	}
	if c.Options.DefaultStringSize > 0 {
		setting.DefaultStringSize = c.Options.DefaultStringSize
	}

	if c.Options.DisableDatetimePrecision == true {
		setting.DisableDatetimePrecision = true
	} else if c.Options.DefaultDatetimePrecision != 3 {
		setting.DefaultDatetimePrecision = &c.Options.DefaultDatetimePrecision
	}

	db, _ := gorm.Open(mysql.New(setting), &gorm.Config{Logger: c.Logger})

	sqlDB, err := db.DB()
	if err != nil {
		log.Println(err)
	}

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
