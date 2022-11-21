package etcd

import (
	"context"
	"time"

	clientV3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type Etcd struct {
	Client *clientV3.Client
}

func New(addrs []string, logger *zap.Logger) (*Etcd, error) {
	client, err := connect(addrs, logger)
	if err != nil {
		return nil, err
	}

	etcd := &Etcd{
		Client: client,
	}

	return etcd, nil
}

func connect(addrs []string, logger *zap.Logger) (*clientV3.Client, error) {
	etcd, err := clientV3.New(clientV3.Config{
		Endpoints:   addrs,
		DialTimeout: 3 * time.Second,
		Logger:      logger,
	})
	if err != nil {
		return nil, err
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	for _, addr := range addrs {
		_, err = etcd.Status(timeoutCtx, addr)
		if err != nil {
			return nil, err
		}
	}
	return etcd, nil
}
