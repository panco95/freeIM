package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	// log *zap.SugaredLogger
	db *mongo.Database
}

func NewMongo(uri, name string, timeout time.Duration, num uint64) (*Client, error) {
	o := options.Client().ApplyURI(uri).SetMaxPoolSize(num).SetTimeout(timeout)
	cli, err := mongo.Connect(context.TODO(), o)
	if err != nil {
		return nil, err
	}
	client := Client{db: cli.Database(name)}
	return &client, nil
}

func (client *Client) Db() *mongo.Database {
	return client.db
}
