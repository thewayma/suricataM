package st

import (
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/thewayma/suricataM/comm/utils"
	"time"
)

type InfluxDBItem struct {
	Name      string                 // TableName, 从metric中获取
	Tags      map[string]string      // MetricData.Tags
	Field     map[string]interface{} // MetricData.Metric <-> MetricData.Value
	Timestamp int64                  // MetricData.Timestamp
}

func (this *InfluxDBItem) String() string {
	return fmt.Sprintf(
		"<TableName:%s, Tags:%v, Field:%v, TS:%d>",
		this.Name,
		this.Tags,
		this.Field,
		this.Timestamp,
	)
}

type InfluxDBInstance struct {
	db client.Client
	bp client.BatchPoints
}

func (r *InfluxDBInstance) InitDB(address, username, password, db, precision string) error {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     address,
		Username: username,
		Password: password,
	})
	if err != nil {
		return err
	}

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  db,
		Precision: precision,
	})
	if err != nil {
		return err
	}

	r.db = c
	r.bp = bp

	return nil
}

func (r *InfluxDBInstance) DeferClose() {
	r.db.Close()
}

func (r *InfluxDBInstance) AddPoints(itemList []interface{}) error {
	for i := 0; i < len(itemList); i++ {
		dbitem := itemList[i].(*InfluxDBItem)

		ti, err := utils.String2Time(utils.UnixTsFormat(dbitem.Timestamp))
		if err != nil {
			ti = time.Now()
		}

		pt, err := client.NewPoint(
			dbitem.Name,
			dbitem.Tags,
			dbitem.Field,
			ti,
		)
		if err != nil {
			return err
		}

		r.bp.AddPoint(pt)
	}

	return nil
}

func (r *InfluxDBInstance) Write() error {
	err := r.db.Write(r.bp)
	if err != nil {
		return err
	}

	return nil
}
