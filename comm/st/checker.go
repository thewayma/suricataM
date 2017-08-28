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

type StrategiesRequest struct {
	IP       string
	Hostname string
}

type Strategy struct {
	Id         int               `json:"id"`
	Metric     string            `json:"metric"`
	Tags       map[string]string `json:"tags"`
	Func       string            `json:"func"`       // e.g. max(#3) all(#3)
	Operator   string            `json:"operator"`   // e.g. < !=
	RightValue float64           `json:"rightValue"` // critical value
	MaxStep    int               `json:"maxStep"`
	Priority   int               `json:"priority"`
	Note       string            `json:"note"`
}

func (this *Strategy) String() string {
	return fmt.Sprintf(
		"<Id:%d, Metric:%s, Tags:%v, %s%s%s MaxStep:%d, P%d, %s, %v>",
		this.Id,
		this.Metric,
		this.Tags,
		this.Func,
		this.Operator,
		utils.ReadableFloat(this.RightValue),
		this.MaxStep,
		this.Priority,
		this.Note,
	)
}

type HostStrategy struct {
	Hostname   string     `json:"hostname"`
	Strategies []Strategy `json:"strategies"`
}

type StrategiesResponse struct {
	HostStrategies []*HostStrategy `json:"hostStrategies"`
}
