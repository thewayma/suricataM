package st

import (
	"fmt"
	"github.com/thewayma/suricataM/comm/utils"
)

type CheckerItem struct {
	Endpoint  string            `json:"endpoint"`
	Metric    string            `json:"metric"`
	Value     float64           `json:"value"`
	Timestamp int64             `json:"timestamp"`
	Type      string            `json:"judgeType"`
	Tags      map[string]string `json:"tags"`
}

func (this *CheckerItem) String() string {
	return fmt.Sprintf("<Endpoint:%s, Metric:%s, Value:%f, Timestamp:%d, JudgeType:%s Tags:%v>",
		this.Endpoint,
		this.Metric,
		this.Value,
		this.Timestamp,
		this.Type,
		this.Tags)
}

func (this *CheckerItem) PrimaryKey() string {
	return utils.Md5(utils.PK(this.Endpoint, this.Metric, this.Tags))
}

type HistoryData struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}
