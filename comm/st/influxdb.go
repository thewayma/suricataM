package st

import (
	"fmt"
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
