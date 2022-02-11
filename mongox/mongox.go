package mongox

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strings"
	"time"
)

type Config struct {
	Host    string
	Options *options.ClientOptions
	Ctx     context.Context
}

func New(host string) *Config {
	return &Config{
		Host:    host,
		Options: options.Client(),
	}
}

func (c *Config) Connect() (*mongo.Client, error) {
	if c.Ctx == nil {
		c.Ctx, _ = context.WithTimeout(context.Background(), 180*time.Second)
	}

	if strings.TrimSpace(c.Host) != "" {
		url := "mongodb://" + c.Host
		c.Options.ApplyURI(url)
	}

	return mongo.Connect(c.Ctx, c.Options)
}

func (c *Config) SetAuth(authMechanism, username, password string) *Config {
	if authMechanism == "PLAIN" {
		c.Options.SetAuth(options.Credential{
			AuthMechanism: authMechanism,
			Username:      username,
			Password:      password,
		})
	} else if authMechanism == "SCRAM" {
		c.Options.SetAuth(options.Credential{
			Username: username,
			Password: password,
		})
	}

	return c
}

func (c *Config) SetReplicaSet(replicaSet string) *Config {
	if strings.TrimSpace(replicaSet) != "" {
		c.Options.SetReplicaSet(replicaSet)
	}
	return c
}

func (c *Config) SetDirect(d bool) *Config {
	c.Options.SetDirect(d)
	return c
}

func (c *Config) SetContext(ctx context.Context) *Config {
	c.Ctx = ctx
	return c
}

func (c *Config) SetRegistry(r *bsoncodec.Registry) *Config {
	c.Options.SetRegistry(r)
	return c
}

func (c *Config) SetPool(minPoolSize, maxPoolSize, maxConnIdleTime uint64) *Config {
	c.Options.SetMinPoolSize(minPoolSize)
	c.Options.SetMaxPoolSize(maxPoolSize)
	c.Options.SetMaxConnIdleTime(time.Duration(maxConnIdleTime) * time.Second)

	return c
}

func (c *Config) SetPoolMonitor() *Config {

	c.Options.SetPoolMonitor(&event.PoolMonitor{
		Event: func(poolEvent *event.PoolEvent) {
			log.Println(">>>>>>>>>>>>>", poolEvent.Type, poolEvent.ConnectionID)
		},
	})

	return c
}
