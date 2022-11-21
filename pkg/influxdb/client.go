package influxdb

import (
	influxdb "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type Point = write.Point

type Client struct {
	influxClient influxdb.Client
	writer       api.WriteAPIBlocking
	reader       api.QueryAPI
}

func New(
	serverURL, authToken string,
	org, bucket string,
) *Client {
	c := &Client{
		influxClient: influxdb.NewClient(serverURL, authToken),
	}

	c.writer = c.influxClient.WriteAPIBlocking(org, bucket)
	c.reader = c.influxClient.QueryAPI(org)

	return c
}

func (c *Client) Writer() api.WriteAPIBlocking {
	return c.writer
}

func (c *Client) Reader() api.QueryAPI {
	return c.reader
}
