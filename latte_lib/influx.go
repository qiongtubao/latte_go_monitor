package latte_lib

import (
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

type InfluxClient struct {
	Url    string
	Db     string
	client client.Client
}

type InfluxConfig struct {
	Url string `json:"url"`
	Db  string `json:"db"`
}

func (c *InfluxClient) Init() error {
	influxdbClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: c.Url,
	})
	if err != nil {
		return err
	}
	c.client = influxdbClient
	return nil
}

func (c *InfluxClient) Send(name string, tags map[string]string, fields map[string]interface{}) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  c.Db,
		Precision: "ns",
	})
	if err != nil {
		return err
	}
	pt, err := client.NewPoint(name, tags, fields, time.Now())
	if err != nil {
		return err
	}
	bp.AddPoint(pt)

	if err := c.client.Write(bp); err != nil {
		return err
	}
	return nil
}
